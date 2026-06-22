package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
)

func InsertProduct(db *sql.DB, name string, category products.Category, banned, favorite bool) error {

	insertSQL := "INSERT INTO products (name, category, banned, favorite) VALUES (?, ?, ?, ?)"

	_, err := db.Exec(insertSQL, name, category, banned, favorite)
	if err != nil {
		return fmt.Errorf("insert db:%w", err)
	}
	return nil
}
