package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Soujuruya/01_SPEC/internal/config"

	"database/sql"

	_ "github.com/lib/pq"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	migrationsPath := flag.String("path", "migrations", "path to migrations folder")
	command := flag.String("command", "", "migration command: up, down, version")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.DB.DSN())
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}(db)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", *migrationsPath),
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("failed to initialize migration: %v", err)
	}

	switch *command {
	case "up":
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migration up failed: %v", err)
		}
		log.Println("migration up successfully")
	case "down":
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("migration down failed: %v", err)
		}
		log.Println("migration down successfully")
	case "version":
		v, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("migration version failed: %v", err)
		}
		log.Printf("migration version: %d, dirty: %v", v, dirty)
	default:
		log.Fatalf("unknown command: %s", *command)
	}
}
