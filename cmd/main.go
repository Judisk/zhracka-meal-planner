package main

import (
	"fmt"
	"log"

	d "foods/internal/dayone"
	f "foods/internal/foodgenerator"
	p "foods/internal/products"
	"foods/internal/storage"
)

func main() {
	fmt.Println(f.GenerateDish("Test name #1", p.Grains, p.Proteins, p.Vegetables, nil))

	fmt.Println(d.GenerateMeals(3, p.Grains, p.Proteins, p.Vegetables, nil))

	db, err := storage.NewDB("products.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
