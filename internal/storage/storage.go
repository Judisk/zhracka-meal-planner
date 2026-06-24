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
        id              INTEGER PRIMARY KEY AUTOINCREMENT,
        name            TEXT NOT NULL CHECK (trim(name) <> ''),
        category        TEXT NOT NULL CHECK (category IN ('grain', 'protein', 'vegetable')),
        banned          INTEGER NOT NULL DEFAULT 0 CHECK (banned IN (0, 1)),
        preference      REAL NOT NULL DEFAULT 1.0 CHECK (preference IN (0.5, 1.0, 1.5)),
        selection_score REAL NOT NULL DEFAULT 1.0
);
`)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("sql get result:%w", err)
	}
	return db, nil
}
