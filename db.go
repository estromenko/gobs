package main

import (
	"database/sql"

	_ "github.com/glebarez/go-sqlite"
)

func GetDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "./db.sqlite3")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS metrics (
		  date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  hostname VARCHAR(64) NOT NULL,
		  name VARCHAR(255) NOT NULL,
		  value VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}
