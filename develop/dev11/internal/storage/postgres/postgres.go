package postgres

import (
	"database/sql"
	"fmt"
	"http-events/internal/config"
	"http-events/internal/storage"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.DBConfig) (*Storage, error) {
	const op = "storage.New"

	connInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.HostDB, cfg.PortDB, cfg.UserDB, cfg.PassDB, cfg.DBName)
	db, err := sql.Open("postgres", connInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveEvent(event storage.Event) error {
	const op = "storage.SaveEvent"

	stmt, err := s.db.Prepare(`INSERT INTO events(user_id, date, event)
	 VALUES($1, $2, $3);`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(event.UserId, event.Date.Time, event.Event)
	if pgerr, ok := err.(*pq.Error); ok {
		if pgerr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrEntryExists)
		}
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateEvent(id int, date time.Time, event string) error {
	const op = "storage.UpdatePerson"

	_, err := s.db.Exec(fmt.Sprintf(`UPDATE events 
	SET %s
	WHERE id=%d;`, getUpdateParams(date, event), id))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	const op = "storage.DeleteEvent"

	stmt, err := s.db.Prepare("DELETE FROM events WHERE id=$1;")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetEventsBetween(userId int, first time.Time, last time.Time) ([]storage.Event, error) {
	const op = "storage.GetEventsBetween"

	stmt, err := s.db.Prepare(`SELECT id, user_id, date, event
	FROM events
	WHERE user_id=$1 AND date>=$2 AND date<$3;`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.Query(userId, first, last)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	res := make([]storage.Event, 0, 100)
	for rows.Next() {
		var event storage.Event
		err := rows.Scan(&event.Id, &event.UserId, &event.Date.Time, &event.Event)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		res = append(res, event)
	}

	return res, nil
}

func getUpdateParams(date time.Time, event string) string {
	var b strings.Builder

	if !date.IsZero() {
		dateStr := fmt.Sprintf("%d/%d/%d", date.Month(), date.Day(), date.Year())
		b.WriteString("date=")
		b.WriteString(fmt.Sprintf("'%s'", dateStr))
		if event != "" {
			b.WriteString(", ")
		}
	}
	if event != "" {
		b.WriteString("event=")
		b.WriteString(fmt.Sprintf("'%s'", event))
	}
	return b.String()
}
