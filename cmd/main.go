package main

import (
	"fmt"
	"log"

	"foods/internal/service"
	s "foods/internal/storage"
)

func main() {

	db, err := s.NewDB("products.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := service.SeedDefaultProductsIfEmpty(db); err != nil {
		log.Fatal(err)
	}

	day, err := service.GenerateAndControlDay(db, 3, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(day)

}
