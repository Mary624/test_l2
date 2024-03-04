package storage

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrEntryExists = errors.New("event exists")
)

type Event struct {
	Id     int        `json:"id"`
	UserId int        `json:"user_id"`
	Date   CustomTime `json:"date"`
	Event  string     `json:"event"`
}

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(b []byte) (err error) {
	date, err := time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	t.Time = date
	return
}

func (e *Event) String() string {
	return fmt.Sprintf("id=%d user_id=%d date=%s event=%s", e.Id, e.UserId, e.Date.String(), e.Event)
}
