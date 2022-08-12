package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "b6e87f3c439c25:dd7a6506@(us-cdbr-east-06.cleardb.net)/heroku_c5002e639cb82ea")

	if err != nil {
		log.Fatal(err)
	}

	return db
}
