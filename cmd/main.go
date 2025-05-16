package main

import (
	"database/sql"
	"log"

	"github.com/whitallee/animal-family-backend/cmd/api"
	"github.com/whitallee/animal-family-backend/config"
	"github.com/whitallee/animal-family-backend/db"
)

func main() {
	cfg := db.PostgresConfig{
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBPort,
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		DBName:   config.Envs.DBName,
		SSLMode:  "disable", // or "require" for production
	}
	db, err := db.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB: Successfully connected!")
}
