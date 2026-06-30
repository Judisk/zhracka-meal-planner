package foodgenerator

import (
	"errors"
	"foods/internal/products"
	"math/rand/v2"
	"testing"
)

func TestGetProd_EmptyList(t *testing.T) {
	var array = []products.Product{}
	_, err := getProd(array, nil)
	if err == nil {
		t.Fatal("getProd() expected error, got nil")
	}
}

func TestGetProd_FixedSeed(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 0))

	g := []products.Product{
		products.NewDefaultProduct("rice", products.Grain),
		products.NewDefaultProduct("buckwheat", products.Grain),
		products.NewDefaultProduct("oats", products.Grain),
	}
	expectedName := "rice"
	prod, err := getProd(g, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)

	}
	if prod.Name != expectedName {
		t.Errorf("expected %q with seed 42, got %q", expectedName, prod.Name)
	}
}

func TestGenerateDish_FixedSeed(t *testing.T) {
	var seed uint64 = 42
	rng := rand.New(rand.NewPCG(seed, 0))

	expectedGrain := "rice"
	expectedProtein := "chicken"
	expectedVegetable := "cucumber"
	testName := "Test name"

	g := []products.Product{
		products.NewDefaultProduct("rice", products.Grain),
		products.NewDefaultProduct("buckwheat", products.Grain),
		products.NewDefaultProduct("oats", products.Grain),
	}
	p := []products.Product{
		products.NewDefaultProduct("egg", products.Protein),
		products.NewDefaultProduct("chicken", products.Protein),
	}
	v := []products.Product{
		products.NewDefaultProduct("cucumber", products.Vegetable),
		products.NewDefaultProduct("tomato", products.Vegetable),
	}

	prod, err := GenerateDish(testName, g, p, v, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prod.Name != testName {
		t.Errorf("expected %q with seed %d, got %q", testName, seed, prod.Name)
	}
	if prod.Grain.Name != expectedGrain {
		t.Errorf("expected %q with seed %d, got %q", expectedGrain, seed, prod.Grain.Name)
	}
	if prod.Protein.Name != expectedProtein {
		t.Errorf("expected %q with seed %d, got %q", expectedProtein, seed, prod.Protein.Name)
	}
	if prod.Vegetable.Name != expectedVegetable {
		t.Errorf("expected %q with seed %d, got %q", expectedVegetable, seed, prod.Vegetable.Name)
	}

}

func TestGeneratorDish_EmptyList(t *testing.T) {
	var seed uint64 = 42
	rng := rand.New(rand.NewPCG(seed, 0))
	testName := "Test name"
	g, p, v := []products.Product{}, []products.Product{}, []products.Product{}
	_, err := GenerateDish(testName, g, p, v, rng)
	if err == nil {
		t.Fatalf("expected an error, got %v", err)
	}
	if !errors.Is(err, ErrEmptyProdList) {
		t.Fatalf("expected an error %v, got %v", ErrEmptyProdList, err)
	}
}
