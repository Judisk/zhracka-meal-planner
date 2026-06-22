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
	var idx int
	if rng != nil {

		idx = rng.IntN(len(array))
	} else {

		idx = rand.IntN(len(array))
	}

	return array[idx], nil
}
