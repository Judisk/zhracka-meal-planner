package main

import (
	"log"
	"math/rand/v2"
	"time"

	"foods/internal/gui"
	"foods/internal/service"
	s "foods/internal/storage"
)

func main() {
	rng := rand.New(rand.NewPCG(
		uint64(time.Now().UnixNano()),
		uint64(time.Now().UnixNano()>>32),
	))

	db, err := s.NewDB("products.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := service.SeedDefaultProductsIfEmpty(db); err != nil {
		log.Fatal(err)
	}

	gui.Run(db, rng)
}
