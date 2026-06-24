package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
)

func InsertProduct(db *sql.DB, name string, category products.Category, banned bool, preference products.PreferenceStatus) error {

	insertSQL := "INSERT INTO products (name, category, banned, preference) VALUES (?, ?, ?, ?)"

	_, err := db.Exec(insertSQL, name, category, banned, preference)
	if err != nil {
		return fmt.Errorf("insert db:%w", err)
	}
	return nil
}
