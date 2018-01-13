package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Open(dataSourceName string) *sql.DB {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	ensureTableExists(db)
	return db
}

func ensureTableExists(db *sql.DB) {
	if _, err := db.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS expenses
(
id INTEGER NOT NULL PRIMARY KEY,
payer TEXT NOT NULL,
amount REAL NOT NULL,
note TEXT NOT NULL DEFAULT '',
date datetime NOT NULL
)`
