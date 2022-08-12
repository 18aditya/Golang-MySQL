package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func Connect() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	User := os.Getenv("DBUSER")
	Pass := os.Getenv("DBPASS")
	Server := os.Getenv("DBSERVER")
	Database := os.Getenv("DB")
	dsm := fmt.Sprintf("%s:%s@(%s)/%s", User, Pass, Server, Database)

	db, err := sql.Open("mysql", dsm)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
