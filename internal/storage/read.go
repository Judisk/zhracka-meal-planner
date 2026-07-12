package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"strings"
)

type queryParams struct {
	query     string
	args      []any
	errString string
}

func SelectAll(db *sql.DB) ([]products.Product, error) {
	q := queryParams{
		query:     "SELECT id, name, category, banned, preference,selection_score FROM products ORDER BY id",
		errString: "query execution",
	}

	return selectAllHelper(db, q)
}

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

	query := "SELECT id, name, category, banned, preference, selection_score FROM products WHERE id = ?"

	var p products.Product

	err := db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Category,
		&p.Banned,
		&p.Preference,
		&p.SelectionScore,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return products.Product{}, fmt.Errorf("product not found: %w", err)
		}
		return products.Product{}, fmt.Errorf("query row: %w", err)
	}

	return p, nil
}

func SelectAllProductsByCategory(db *sql.DB, category products.Category) ([]products.Product, error) {

	q := queryParams{
		query: `
		SELECT id, name, category, banned, preference, selection_score
		FROM products
		WHERE category = ?
		ORDER BY id
	`,
		args:      []any{category},
		errString: "query execution",
	}
	return selectAllHelper(db, q)
}

func SelectReadyProductsByCategory(db *sql.DB, category products.Category) ([]products.Product, error) {
	q := queryParams{
		query: `
		SELECT id, name, category, banned, preference, selection_score
		FROM products
		WHERE category = ? AND banned = 0 AND selection_score > 0
		ORDER BY id
	`,
		args:      []any{category},
		errString: "query ready products",
	}
	return selectAllHelper(db, q)
}

func SelectUnbannedProductsByCategory(db *sql.DB, category products.Category) ([]products.Product, error) {

	q := queryParams{
		query: `
		SELECT id, name, category, banned, preference, selection_score
		FROM products
		WHERE category = ? AND banned = 0
		ORDER BY id
	`,
		args:      []any{category},
		errString: "query allowed products",
	}
	return selectAllHelper(db, q)
}

func SelectAllFiltered(db *sql.DB, category *products.Category, banned *bool, preference *products.PreferenceStatus) ([]products.Product, error) {
	conditions := []string{}
	args := []any{}

	if category != nil {
		conditions = append(conditions, "category = ?")
		args = append(args, *category)
	}
	if banned != nil {
		conditions = append(conditions, "banned = ?")
		args = append(args, *banned)
	}
	if preference != nil {
		conditions = append(conditions, "preference = ?")
		args = append(args, *preference)
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	q := queryParams{
		query:     "SELECT id, name, category, banned, preference, selection_score FROM products " + where + " ORDER BY id",
		args:      args,
		errString: "query filtered products",
	}
	return selectAllHelper(db, q)
}

func selectAllHelper(db *sql.DB, q queryParams) ([]products.Product, error) {
	rows, err := db.Query(q.query, q.args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", q.errString, err)
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
			&p.Preference,
			&p.SelectionScore,
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
