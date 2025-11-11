// Package db contains a tiny wrapper around an SQLite database used to
// persist scheduler tasks. It provides initialization, a package-global DB
// handle and helper functions for CRUD operations located in other files.
package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const (
	schema = `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		-- date stored as YYYYMMDD (e.g. 20060102). enforce 8 numeric chars.
		date CHAR(8) NOT NULL CHECK(date GLOB '[0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]'),
		title VARCHAR NOT NULL,
		comment TEXT,
		-- repeat rules string, limited to 128 characters
		repeat VARCHAR CHECK(length(repeat) <= 128)
	);
	-- index on date to speed up lookups by date
	CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
	`
)

// package-global database handle (keeps one open connection)
var db *sql.DB

// DB returns the package-global database handle. It may be nil if Init wasn't called.
func DB() *sql.DB { return db }

// Init opens (or creates) the SQLite database at dbFile, assigns the handle to
// the package-global variable and creates the schema if the file did not exist.
// It is safe to call once during application startup.
func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	// Open and keep the database connection in package-global variable `db`.
	d, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	// create schema if needed
	if install {
		if _, err = db.Exec(schema); err != nil {
			db.Close()
			db = nil
			return err
		}
	}

	// assign to package-global
	db = d

	return nil
}

// Close closes the package-global database connection if open and clears the
// package-global variable. Call this during graceful shutdown.
func Close() error {
	if db == nil {
		return nil
	}
	err := db.Close()
	db = nil
	return err
}
