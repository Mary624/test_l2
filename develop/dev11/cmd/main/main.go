package main

import (
	"http-events/internal/config"
	"http-events/internal/server"
	"http-events/internal/storage/postgres"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullPath := filepath.Join(filepath.Join(path, "../.."), ".env")
	err = godotenv.Load(fullPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cfg := config.MustLoad()

	db, err := postgres.New(cfg.DBConfig)
	if err != nil {
		panic("can't connect to db")
	}
	srv := server.New(cfg, db, db)
	srv.Run(cfg.Port)
}
