package service

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"foods/internal/storage"
)

func SeedDefaultProductsIfEmpty(db *sql.DB) error {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return fmt.Errorf("count products: %w", err)
	}

	if count == 0 {
		return insertDefaultProducts(db)
	}

	return nil
}

func insertDefaultProducts(db *sql.DB) error {

	defaults := []products.Product{
		products.NewDefaultProduct("рис", products.Grain),
		products.NewDefaultProduct("овес", products.Grain),
		products.NewDefaultProduct("гречка", products.Grain),

		products.NewDefaultProduct("яйцо", products.Protein),
		products.NewDefaultProduct("курица", products.Protein),

		products.NewDefaultProduct("томат", products.Vegetable),
		products.NewDefaultProduct("огурец", products.Vegetable),
	}

	for _, p := range defaults {
		if err := storage.InsertProduct(db, p); err != nil {
			return fmt.Errorf("insert default product %q: %w", p.Name, err)
		}
	}

	return nil
}

