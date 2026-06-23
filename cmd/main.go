package main

import (
	"database/sql"
	"fmt"
	"log"

	"foods/internal/dayone"
	"foods/internal/foodgenerator"
	"foods/internal/products"
	"foods/internal/storage"
)

func main() {

	db, err := storage.NewDB("products.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := seedDefaultProductsIfEmpty(db); err != nil {
		log.Fatal(err)
	}
	grains, err := storage.SelectProductsByCategory(db, products.Grain)
	if err != nil {
		log.Fatal(err)
	}
	proteins, err := storage.SelectProductsByCategory(db, products.Protein)
	if err != nil {
		log.Fatal(err)
	}
	vegetables, err := storage.SelectProductsByCategory(db, products.Vegetable)
	if err != nil {
		log.Fatal(err)
	}

	day, err := dayone.GenerateMeals(3, grains, proteins, vegetables, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(day)

	dish, err := foodgenerator.GenerateDish("Test name #1", grains, proteins, vegetables, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dish)
}

func seedDefaultProductsIfEmpty(db *sql.DB) error {
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
		{Name: "рис", Category: products.Grain},
		{Name: "гречка", Category: products.Grain},
		{Name: "овес", Category: products.Grain},

		{Name: "яйцо", Category: products.Protein},
		{Name: "курица", Category: products.Protein},

		{Name: "томат", Category: products.Vegetable},
		{Name: "огурец", Category: products.Vegetable},
	}

	for _, p := range defaults {
		if err := storage.InsertProduct(db, p.Name, p.Category, p.Banned, p.Favorite); err != nil {
			return fmt.Errorf("insert default product %q: %w", p.Name, err)
		}
	}

	return nil
}
