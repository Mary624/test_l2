package server

import (
	"fmt"
	"http-events/internal/config"
	"http-events/internal/handlers/get"
	"http-events/internal/handlers/post"
	"http-events/internal/middleware"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	mux *http.ServeMux
	log *slog.Logger
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func New(cfg config.Config, change post.ChangeEvents, getter get.Getter) *Server {
	mux := http.NewServeMux()
	log := setupLogger(cfg.Env)

	// POST
	create := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { post.Post(change, post.Create, w, r) })
	update := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { post.Post(change, post.Update, w, r) })
	delete := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { post.Post(change, post.Delete, w, r) })

	mux.Handle("/create_event", middleware.MiddlewareLog(log, create))
	mux.Handle("/update_event", middleware.MiddlewareLog(log, update))
	mux.Handle("/delete_event", middleware.MiddlewareLog(log, delete))

	//GET
	getDay := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { get.GetByDate(getter, get.Day, w, r) })
	getWeek := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { get.GetByDate(getter, get.Week, w, r) })
	getMonth := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { get.GetByDate(getter, get.Month, w, r) })

	mux.Handle("/events_for_day", middleware.MiddlewareLog(log, getDay))
	mux.Handle("/events_for_week", middleware.MiddlewareLog(log, getWeek))
	mux.Handle("/events_for_month", middleware.MiddlewareLog(log, getMonth))

	return &Server{
		mux: mux,
		log: log,
	}
}

func (s *Server) Run(port int) error {
	s.log.Info("start server")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), s.mux)
	return err
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
