package main

import (
	"database/sql"
	"log"

	"github.com/whitallee/animal-family-backend/cmd/api"
	"github.com/whitallee/animal-family-backend/config"
	"github.com/whitallee/animal-family-backend/db"
)

//	@title			Animal Family API
//	@version		2.0
//	@description	REST API for managing animals, enclosures and care tasks.
//	@description
//	@description	Only v2 routes are documented here. v1 remains served at /api/v1 for
//	@description	backwards compatibility but is not part of this contract.

//	@BasePath	/api/v2

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Raw JWT token from POST /users/login. Sent verbatim with no "Bearer " prefix.
func main() {
	cfg := db.PostgresConfig{
		Host:     config.Envs.DBHost,
		Port:     config.Envs.DBPort,
		User:     config.Envs.DBUser,
		Password: config.Envs.DBPassword,
		DBName:   config.Envs.DBName,
		SSLMode:  config.Envs.DBSSLMode,
	}
	db, err := db.NewPostgresStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(":"+config.Envs.Port, db)
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
