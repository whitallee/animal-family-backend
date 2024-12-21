package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN()) //ATTEMPTING TO CONNECT WITHOUT DSN FORMAT
	//db, err := sql.Open("mysql", "mysql://root:oqtTuMGuqlQSFatEpiicJJcdGfWbANtd@junction.proxy.rlwy.net:58595/railway")
	if err != nil {
		println("error in db.go NewMySqlStorage method")
		log.Fatal(err)
	}

	return db, nil
}
