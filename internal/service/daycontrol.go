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
	if n > 6 {
		return dayone.Day{}, fmt.Errorf("generate and control day: %d is too many meals for one day", n)
	}
	if n <= 0 {
		return dayone.Day{}, fmt.Errorf("generate and control day: %d must be positive", n)
	}

	grains, err := storage.SelectReadyProductsByCategory(db, products.Grain)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate and control day: %w", err)
	}
	proteins, err := storage.SelectReadyProductsByCategory(db, products.Protein)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate and control day: %w", err)
	}

	vegetables, err := storage.SelectReadyProductsByCategory(db, products.Vegetable)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate and control day: %w", err)
	}

	banned, err := storage.SelectBannedProducts(db)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate and control day: %w", err)
	}

	day, err := dayone.GenerateMeals(n, grains, proteins, vegetables, banned, rng)
	if err != nil {
		return dayone.Day{}, fmt.Errorf("generate day: %w", err)
	}
	arr := []products.Product{}
	for _, elem := range day.Meals {
		arr = append(arr, elem.Grain, elem.Protein, elem.Vegetable)
	}

	if err := storage.ManyResets(db, arr...); err != nil {
		return dayone.Day{}, fmt.Errorf("reset scores: %w", err)
	}
	if err := storage.UpdateSelectionScore(db); err != nil {
		return dayone.Day{}, fmt.Errorf("update scores: %w", err)
	}
	return day, nil
}
