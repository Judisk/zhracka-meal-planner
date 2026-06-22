package storage

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("sql open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sql ping: %w", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
    id       INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK (trim(name) <> ''),
    category TEXT NOT NULL CHECK (category IN ('grain', 'protein', 'vegetable')),
    banned   BOOLEAN DEFAULT FALSE,
    favorite BOOLEAN DEFAULT FALSE
);
	`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("sql get result:%w", err)
	}
	return db, nil
}

