package service

import (
	"database/sql"
	"fmt"
	"foods/internal/dayone"
	"foods/internal/products"
	"foods/internal/storage"
	"math/rand/v2"
)

func GenerateAndControlDay(db *sql.DB, n int, rng *rand.Rand) (dayone.Day, error) {
	grains, err := storage.SelectAllowedProductsByCategory(db, products.Grain)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate and control day: %w", err)
	}
	proteins, err := storage.SelectAllowedProductsByCategory(db, products.Protein)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate and control day: %w", err)
	}

	vegetables, err := storage.SelectAllowedProductsByCategory(db, products.Vegetable)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate and control day: %w", err)
	}

	day, err := dayone.GenerateMeals(n, grains, proteins, vegetables, rng)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate day: %w", err)
	}
	arr := []products.Product{}
	for _, elem := range day.Meals {
		arr = append(arr, elem.Grain, elem.Protein, elem.Vegetable)
	}
	if err := storage.ManyResets(db, arr...); err != nil {
		return dayone.Day{}, fmt.Errorf("controle day: %w", err)
	}
	return day, nil
}
