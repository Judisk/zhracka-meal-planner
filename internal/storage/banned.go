package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"strings"
)

func SelectBannedProducts(db *sql.DB) (products.BlockedProducts, error) {
	query := "SELECT id, name FROM products WHERE banned = 1 ORDER BY id"

	rows, err := db.Query(query)
	if err != nil {
		return products.BlockedProducts{}, fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

	blocked := products.BlockedProducts{
		ByID:   map[products.ProductID]bool{},
		ByName: map[string]bool{},
	}

	for rows.Next() {
		var id products.ProductID
		var name string

		if err := rows.Scan(
			&id,
			&name,
		); err != nil {
			return products.BlockedProducts{}, fmt.Errorf("scan row: %w", err)
		}

		name = normalizeProductName(name)

		blocked.ByID[id] = true
		blocked.ByName[name] = true

	}

	if err = rows.Err(); err != nil {
		return products.BlockedProducts{}, fmt.Errorf("row iteration: %w", err)
	}

	return blocked, nil

}

func normalizeProductName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
