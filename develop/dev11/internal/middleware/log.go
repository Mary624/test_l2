package middleware

import (
	"bytes"
	"encoding/json"
	"http-events/internal/handlers"
	logresponsewriter "http-events/internal/handlers/log_response_writer"
	"http-events/internal/storage"
	"io"
	"log/slog"
	"net/http"
)

func MiddlewareLog(log *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug(r.Method, slog.String("Endpoint", r.URL.Path))

		if r.Method == http.MethodPost {

			body, err := io.ReadAll(r.Body)
			if err == nil {
				var event storage.Event
				r.Body = io.NopCloser(bytes.NewBuffer(body))
				err = json.Unmarshal(body, &event)
				if err == nil {
					log.Debug(event.String())
				}
			}
		}
		lrw := logresponsewriter.NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		if lrw.StatusCode >= 500 {
			var errMes handlers.ErrorMessage
			err := json.Unmarshal(lrw.Data, &errMes)
			if err == nil {
				log.Error(errMes.Error)
			}
		}
	})
}
