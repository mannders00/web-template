package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		log.Fatal(err)
	}

	stmt := `CREATE TABLE users (email TEXT UNIQUE, hashed_password TEXT)`
	_, err = db.Exec(stmt)
	if err != nil {
		log.Print(err)
	}
}
