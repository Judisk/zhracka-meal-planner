package foodgenerator

import (
	"fmt"
	"foods/internal/products"
	"math/rand/v2"
)

type Dish struct {
	Name      string
	Grain     products.Product
	Protein   products.Product
	Vegetable products.Product
}

func GenerateDish(name string, grains, proteins, vegetables []products.Product, rng *rand.Rand) (Dish, error) {
	g, err := getProd(grains, rng)
	if err != nil {
		return Dish{}, fmt.Errorf("pick grain: %w", err)
	}
	p, err := getProd(proteins, rng)
	if err != nil {
		return Dish{}, fmt.Errorf("pick protein: %w", err)
	}
	v, err := getProd(vegetables, rng)
	if err != nil {
		return Dish{}, fmt.Errorf("pick vegetable: %w", err)
	}
	return Dish{
		Name:      name,
		Grain:     g,
		Protein:   p,
		Vegetable: v,
	}, nil
}

func getProd(array []products.Product, rng *rand.Rand) (products.Product, error) {
	if len(array) == 0 {
		return products.Product{}, fmt.Errorf("getProd: empty product list")
	}
	var total float64
	for _, p := range array {
		total += p.SelectionScore
	}
	r := rng.Float64() * total
	for _, p := range array {
		r -= p.SelectionScore
		if r <= 0 {
			return p, nil
		}
	}
	return products.Product{}, fmt.Errorf("getProd: random choice did not select a product")
}
