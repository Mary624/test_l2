package post

import (
	"encoding/json"
	"http-events/internal/handlers"
	"http-events/internal/storage"
	"net/http"
	"strings"
	"time"
)

type ChangeEvents interface {
	SaveEvent(storage.Event) error
	UpdateEvent(int, time.Time, string) error
	DeleteEvent(int) error
}

const (
	Create = iota
	Update
	Delete
)

func Post(changeEvents ChangeEvents, action int, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	var event storage.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		handlers.WriteError(w, err.Error(), http.StatusBadGateway)
		return
	}

	switch action {
	case Create:
		if strings.Contains(event.Event, ";") || event.UserId < 1 {
			handlers.WriteError(w, handlers.ErrWrongInput.Error(), http.StatusBadRequest)
			return
		}
		err = changeEvents.SaveEvent(event)
	case Update:
		if event.Date.IsZero() && event.Event == "" || strings.Contains(event.Event, ";") {
			handlers.WriteError(w, handlers.ErrWrongInput.Error(), http.StatusBadRequest)
			return
		}
		err = changeEvents.UpdateEvent(event.Id, event.Date.Time, event.Event)
	case Delete:
		if event.Id < 1 {
			handlers.WriteError(w, handlers.ErrWrongInput.Error(), http.StatusBadRequest)
			return
		}
		err = changeEvents.DeleteEvent(event.Id)
	default:
		handlers.WriteError(w, "wrong period", http.StatusBadGateway)
		return
	}
	if err != nil {
		handlers.WriteError(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`"result":"done"`))
}
