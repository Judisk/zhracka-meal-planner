package main

import (
	"fmt"

	d "foods/internal/dayone"
	f "foods/internal/foodgenerator"
	p "foods/internal/products"
)

func main() {
	fmt.Println(f.GenerateDish("Test name #1", p.Grains, p.Proteins, p.Vegetables, nil))

	fmt.Println(d.GenerateMeals(3, p.Grains, p.Proteins, p.Vegetables, nil))
}
