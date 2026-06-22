package storage

import (
	"database/sql"
	"fmt"
)

func DeleteProductsByName(db *sql.DB, name string) error {
	deleteSQL := "DELETE FROM products WHERE name = ?"

	_, err := db.Exec(deleteSQL, name)
	if err != nil {
		return fmt.Errorf("delete products by name: %w", err)
	}
	return nil
}
