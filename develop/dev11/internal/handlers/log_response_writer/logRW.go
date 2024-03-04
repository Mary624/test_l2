package logresponsewriter

import "net/http"

type loggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Data       []byte
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, nil}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(data []byte) (int, error) {
	lrw.Data = data
	return lrw.ResponseWriter.Write(data)
}
