package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
)

func SelectProductIDsByCategory(db *sql.DB, category products.Category) ([]products.ProductID, error) {
	query := "SELECT id FROM products WHERE category = ? ORDER BY id"

	rows, err := db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

	var ids []products.ProductID

	for rows.Next() {
		var id products.ProductID

		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration: %w", err)
	}

	return ids, nil
}

func SelectProductByID(db *sql.DB, id products.ProductID) (products.Product, error) {

	query := "SELECT id, name, category, banned, favorite FROM products WHERE id = ?"

	var p products.Product

	err := db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Category,
		&p.Banned,
		&p.Favorite,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return products.Product{}, fmt.Errorf("product not found: %w", err)
		}
		return products.Product{}, fmt.Errorf("query row: %w", err)
	}

	return p, nil
}

func SelectProductsByCategory(db *sql.DB, category products.Category) ([]products.Product, error) {
	query := `
		SELECT id, name, category, banned, favorite
		FROM products
		WHERE category = ?
		ORDER BY id
	`

	rows, err := db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("query execution: %w", err)
	}

	defer rows.Close()

	var result []products.Product

	for rows.Next() {
		var p products.Product

		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Category,
			&p.Banned,
			&p.Favorite,
		); err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		result = append(result, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration: %w", err)
	}

	return result, nil
}
