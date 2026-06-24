package main

import (
	"database/sql"
	"fmt"
	"log"

	d "foods/internal/dayone"
	f "foods/internal/foodgenerator"
	p "foods/internal/products"
	s "foods/internal/storage"
)

func main() {

	db, err := s.NewDB("products.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := seedDefaultProductsIfEmpty(db); err != nil {
		log.Fatal(err)
	}
	grains, err := s.SelectAllowedProductsByCategory(db, p.Grain)
	if err != nil {
		log.Fatal(err)
	}
	proteins, err := s.SelectAllowedProductsByCategory(db, p.Protein)
	if err != nil {
		log.Fatal(err)
	}
	vegetables, err := s.SelectAllowedProductsByCategory(db, p.Vegetable)
	if err != nil {
		log.Fatal(err)
	}

	day, err := d.GenerateMeals(3, grains, proteins, vegetables, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(day)

	dish, err := f.GenerateDish("Test name #1", grains, proteins, vegetables, nil)
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

	defaults := []p.Product{
		p.NewProduct("рис", p.Grain),
		p.NewProduct("овес", p.Grain),
		p.NewProduct("гречка", p.Grain),

		p.NewProduct("яйцо", p.Protein),
		p.NewProduct("курица", p.Protein),

		p.NewProduct("томат", p.Vegetable),
		p.NewProduct("огурец", p.Vegetable),
	}

	for _, p := range defaults {
		if err := s.InsertProduct(db, p.Name, p.Category, p.Banned, p.Preference); err != nil {
			return fmt.Errorf("insert default product %q: %w", p.Name, err)
		}
	}

	return nil
}
