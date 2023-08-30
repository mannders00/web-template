package main

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
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

func register(email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("INSERT INTO users (email, hashed_password) VALUES (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(email, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func login(email string, password string) error {
	var hashedPassword []byte

	err := db.QueryRow("SELECT hashed_password FROM users WHERE email = ?", email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("User not found")
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		return fmt.Errorf("Invalid password")
	}

	return nil
}
