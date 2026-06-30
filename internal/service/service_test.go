package service

import (
	"database/sql"
	"foods/internal/products"
	"foods/internal/storage"
	"math/rand/v2"
	"testing"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := storage.NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return db
}

func setupDBAndRand(t *testing.T) (*sql.DB, *rand.Rand) {
	t.Helper()
	return setupDB(t), rand.New(rand.NewPCG(42, 0))
}
func TestService_NOutOfCount(t *testing.T) {
	tests := []struct {
		name string
		N    int
	}{{
		name: "n exceeds maximum",
		N:    7,
	}, {
		name: "n is below minimum",
		N:    0,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, rng := setupDBAndRand(t)
			defer db.Close()
			_, err := GenerateAndControlDay(db, tt.N, rng)
			if err == nil {
				t.Fatalf("expected an error, got %v", err)
			}
		})

	}
}

// TODO: merge these tests into the table-driven test

func TestService_GenerateAndControlDaySuccess(t *testing.T) {
	db, rng := setupDBAndRand(t)
	defer db.Close()
	N := 3

	usedProds := make(map[string]bool)

	prods := []products.Product{
		products.NewDefaultProduct("rice", products.Grain),
		products.NewDefaultProduct("buckwheat", products.Grain),
		products.NewDefaultProduct("oats", products.Grain),

		products.NewDefaultProduct("egg", products.Protein),
		products.NewDefaultProduct("chicken", products.Protein),
		products.NewDefaultProduct("ham", products.Protein),

		products.NewDefaultProduct("herbs", products.Vegetable),
		products.NewDefaultProduct("cucumber", products.Vegetable),
		products.NewDefaultProduct("tomato", products.Vegetable),
	}

	for _, p := range prods {
		if err := storage.InsertProduct(db, p); err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}

	day, err := GenerateAndControlDay(db, N, rng)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	for _, elem := range day.Meals {
		if _, ok := usedProds[elem.Grain.Name]; ok {
			t.Errorf("unexpected repeat %q", elem.Grain.Name)
		}

		if _, ok := usedProds[elem.Vegetable.Name]; ok {
			t.Errorf("unexpected repeat %q", elem.Vegetable.Name)
		}

		if _, ok := usedProds[elem.Protein.Name]; ok {
			t.Errorf("unexpected repeat %q", elem.Protein.Name)
		}
		usedProds[elem.Grain.Name] = true
		usedProds[elem.Vegetable.Name] = true
		usedProds[elem.Protein.Name] = true
	}

}

func TestService_GenerateAndControlClosedDB(t *testing.T) {
	db, rng := setupDBAndRand(t)
	N := 3

	db.Close()

	_, err := GenerateAndControlDay(db, N, rng)
	if err == nil {
		t.Fatalf("expected an error got %v", err)
	}

}

func TestService_GenerateAndControlDayEmptyDB(t *testing.T) {
	db, rng := setupDBAndRand(t)
	defer db.Close()
	N := 3

	_, err := GenerateAndControlDay(db, N, rng)
	if err == nil {
		t.Fatalf("expected an error got %v", err)
	}

}

func TestService_GenerateAndControlDayNMoreThanProds(t *testing.T) {
	db, rng := setupDBAndRand(t)
	defer db.Close()
	N := 3

	prods := []products.Product{
		products.NewDefaultProduct("rice", products.Grain),
		products.NewDefaultProduct("buckwheat", products.Grain),

		products.NewDefaultProduct("egg", products.Protein),
		products.NewDefaultProduct("chicken", products.Protein),
		products.NewDefaultProduct("ham", products.Protein),

		products.NewDefaultProduct("herbs", products.Vegetable),
		products.NewDefaultProduct("cucumber", products.Vegetable),
		products.NewDefaultProduct("tomato", products.Vegetable),
	}

	for _, p := range prods {
		if err := storage.InsertProduct(db, p); err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}

	_, err := GenerateAndControlDay(db, N, rng)
	if err == nil {
		t.Fatalf("expected an error got %v", err)
	}

}
