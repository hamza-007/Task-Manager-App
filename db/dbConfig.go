package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Connect() *sql.DB {
	db, e := sql.Open("sqlite3", "./db/data.db")
	if e != nil {
		log.Fatalf("Error: %v", e)
		return nil
	}

	if e := db.Ping(); e != nil {
		log.Fatalf("Error: %v", e)
		return nil
	}
	return db
}
