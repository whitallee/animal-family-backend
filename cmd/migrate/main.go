package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/whitallee/animal-family-backend/config"
	"github.com/whitallee/animal-family-backend/db"
)

func main() {
	db, err := db.NewPostgresStorage(db.PostgresConfig{
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBPort,
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		DBName:   config.Envs.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		println("error in main.go NewPostgresStorage method")
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]

	if cmd == "up" {
		if err := m.Up(); err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
