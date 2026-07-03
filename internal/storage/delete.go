package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
)

func DeleteProductsByName(db *sql.DB, name string) error {
	deleteSQL := "DELETE FROM products WHERE name = ?"

	_, err := db.Exec(deleteSQL, name)
	if err != nil {
		return fmt.Errorf("delete products by name: %w", err)
	}
	return nil
}

func DeleteProductsByID(db *sql.DB, id products.ProductID) error {
	deleteSQL := "DELETE FROM products WHERE id = ?"

	_, err := db.Exec(deleteSQL, id)
	if err != nil {
		return fmt.Errorf("delete products by id: %w", err)
	}
	return nil
}
