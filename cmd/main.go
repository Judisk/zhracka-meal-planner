package main

import (
	"log"
	"math/rand/v2"
	"os"
	"time"

	"foods/internal/gui"
	"foods/internal/service"
	s "foods/internal/storage"
)

func main() {
	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
func run() error {
	rng := rand.New(rand.NewPCG(
		uint64(time.Now().UnixNano()),
		uint64(time.Now().UnixNano()>>32),
	))

	db, err := s.NewDB("products.db")
	if err != nil {
		return err
	}

	defer db.Close()

	if err := service.SeedDefaultProductsIfEmpty(db); err != nil {
		return err
	}

	return gui.Run(db, rng)

}
