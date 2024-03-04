package handlers

import (
	"encoding/json"
	"errors"
	"http-events/internal/storage"
	"net/http"
)

var ErrWrongInput = errors.New("invalid input data")

type ErrorMessage struct {
	Error string `json:"error"`
}

type ResultMessage struct {
	Result []storage.Event `json:"result"`
}

func WriteError(w http.ResponseWriter, mes string, status int) {
	w.WriteHeader(status)
	var errMes ErrorMessage
	errMes.Error = mes
	b, err := json.Marshal(errMes)
	if err == nil {
		w.Write(b)
	}
}
