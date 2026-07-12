package storage

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
)

func UpdateSelectionScore(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	updateSQLReady := `
	UPDATE products
	SET selection_score = selection_score + preference
	WHERE banned = 0 
		AND selection_score >=1
		AND selection_score < 12
	`

	updateSQLCooldown := `
	UPDATE products
	SET selection_score = selection_score + 1.0
	WHERE banned = 0 AND selection_score IN (0,-1)
	`

	if _, err := tx.Exec(updateSQLReady); err != nil {
		return fmt.Errorf("update ready product scores: %w", err)
	}
	if _, err := tx.Exec(updateSQLCooldown); err != nil {
		return fmt.Errorf("update cooldown product scores: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit selection score update: %w", err)
	}
	return nil
}

func ResetProductScore(tx *sql.Tx, p products.Product) error {
	updateSQLCD := `
	UPDATE products
	SET selection_score = -1.0
	WHERE banned =0 AND id = ?
	`
	_, err := tx.Exec(updateSQLCD, p.ID)
	if err != nil {
		return fmt.Errorf("reset cooldown by id: %w", err)
	}
	return nil
}

func ManyResets(db *sql.DB, prods ...products.Product) error {
	if len(prods) == 0 {
		return fmt.Errorf("many resets: empty prods")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	for _, elem := range prods {
		if err := ResetProductScore(tx, elem); err != nil {
			return fmt.Errorf("reset many products: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit selection score update: %w", err)
	}
	return nil
}

func UpdateProductInfo(db *sql.DB, p products.Product) error {
	query := "UPDATE products SET name=?, category=?, banned=?, preference=? WHERE id=?"
	_, err := db.Exec(query, p.Name, p.Category, p.Banned, p.Preference, p.ID)
	if err != nil {
		return fmt.Errorf("update product info: %w", err)
	}
	return nil
}
