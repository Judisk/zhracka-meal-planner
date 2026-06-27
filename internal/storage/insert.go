package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
)

func InsertProduct(db *sql.DB, p products.Product) error {

	insertSQL := "INSERT INTO products (name, category, banned, preference) VALUES (?, ?, ?, ?)"

	_, err := db.Exec(insertSQL, p.Name, p.Category, p.Banned, p.Preference)
	if err != nil {
		return fmt.Errorf("insert db:%w", err)
	}
	return nil
}
