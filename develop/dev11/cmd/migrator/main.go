package main

import (
	"errors"
	"fmt"
	"http-events/internal/config"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fullPath := filepath.Join(path, ".env")
	err = godotenv.Load(fullPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	cfg := config.MustLoadDB()
	m, err := migrate.New(
		"file://"+cfg.MigrationsPath,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?x-migrations-table=%s&sslmode=disable", cfg.UserDB, cfg.PassDB, cfg.HostDB, cfg.PortDB, cfg.DBName, cfg.MigrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")

			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}
