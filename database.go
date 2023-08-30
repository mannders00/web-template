package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func initDB() {
	db, err := sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt := `CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, hashed_password TEXT)`
	_, err = db.Exec(stmt)
	if err != nil {
		log.Fatal(err)
	}
}
