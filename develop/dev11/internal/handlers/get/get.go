package get

import (
	"encoding/json"
	"fmt"
	"http-events/internal/handlers"
	"http-events/internal/storage"
	"net/http"
	"sort"
	"strconv"
	"time"
)

type Getter interface {
	GetEventsBetween(int, time.Time, time.Time) ([]storage.Event, error)
}

const (
	Day = iota
	Week
	Month
)

func get(w http.ResponseWriter, r *http.Request) (int, time.Time, error) {
	userId := r.URL.Query().Get("user_id")
	date := r.URL.Query().Get("date")
	if date == "" || userId == "" {
		return 0, time.Time{}, fmt.Errorf("empty params")
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("can't get user id")
	}

	dateTime, err := time.Parse("2006-01-02 15:04:05", string(date)+" 00:00:00")
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("can't get date")
	}

	return userIdInt, dateTime, nil
}

func GetByDate(getter Getter, period int, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	user_id, date, err := get(w, r)
	if err != nil {
		handlers.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var first, last time.Time
	switch period {
	case Day:
		first, last = getDay(date)
	case Week:
		first, last = getWeek(date)
	case Month:
		first, last = getMonth(date)
	default:
		handlers.WriteError(w, "wrong period", http.StatusBadGateway)
		return
	}

	res, err := getter.GetEventsBetween(user_id, first, last)
	if err != nil {
		handlers.WriteError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].UserId < res[j].UserId {
			return true
		}
		if res[i].Date.Time.Before(res[j].Date.Time) {
			return true
		}
		return res[i].Id < res[j].Id
	})

	var resMes handlers.ResultMessage
	resMes.Result = res
	b, err := json.Marshal(resMes)
	if err != nil {
		handlers.WriteError(w, err.Error(), http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func getDay(date time.Time) (time.Time, time.Time) {
	var first time.Time = date
	last := first.Add(time.Hour * 24)
	return first, last
}

func getWeek(date time.Time) (time.Time, time.Time) {
	weekday := date.Weekday()
	weekday--
	if weekday == -1 {
		weekday = 6
	}
	monday := date.Add(-time.Duration(weekday) * time.Hour * 24)
	sunday := monday.Add(time.Hour * 24 * 7)

	return monday, sunday
}

func getMonth(date time.Time) (time.Time, time.Time) {
	day := date.Day()
	first := date.Add(-time.Hour * 24 * time.Duration(day-1))
	last := first.Add(time.Duration(maxDayByMonth(date)) * time.Hour * 24)

	return first, last
}

func maxDayByMonth(date time.Time) int {
	switch date.Month() {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.April, time.June, time.September, time.November:
		return 30
	case time.February:
		if isLeapYear(date) {
			return 29
		}
		return 28
	}
	return 0
}

func isLeapYear(date time.Time) bool {
	if date.Year()%400 == 0 {
		return true
	}
	return date.Year()%4 == 0 && date.Year()%100 != 0
}
