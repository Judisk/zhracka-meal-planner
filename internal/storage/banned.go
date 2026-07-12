package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
)

func SelectBannedProducts(db *sql.DB) (products.BlockedProducts, error) {
	query := "SELECT id, name FROM products WHERE banned = 1 ORDER BY id"

	rows, err := db.Query(query)
	if err != nil {
		return products.BlockedProducts{}, fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

	byID := make(map[products.ProductID]bool)
	byName := []string{}

	for rows.Next() {
		var id products.ProductID
		var name string

		if err := rows.Scan(
			&id,
			&name,
		); err != nil {
			return products.BlockedProducts{}, fmt.Errorf("scan row: %w", err)
		}

		byID[id] = true
		byName = append(byName, name)

	}

	if err = rows.Err(); err != nil {
		return products.BlockedProducts{}, fmt.Errorf("row iteration: %w", err)
	}

	return products.NewBlockedProducts(byID, byName), nil

}
