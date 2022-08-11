package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:admin2021@tcp(127.0.0.1:3306)/adit_schema")

	if err != nil {
		log.Fatal("asdasdasda ? ", err)
	}

	return db
}
