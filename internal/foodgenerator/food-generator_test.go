package foodgenerator

import (
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
		products.NewDefaultProduct("рис", products.Grain),
		products.NewDefaultProduct("гречка", products.Grain),
		products.NewDefaultProduct("овес", products.Grain),
	}
	expectedName := "рис"
	prod, err := getProd(g, rng)
	if err != nil {
		t.Fatalf("ошибка: %v", err)

	}
	if prod.Name != expectedName {
		t.Errorf("ожидали %q при сиде 42, получили %q", expectedName, prod.Name)
	}
}

func TestGenerateDish_FixedSeed(t *testing.T) {
	var seed uint64 = 42
	rng := rand.New(rand.NewPCG(seed, 0))

	expectedGrain := "рис"
	expectedProtein := "курица"
	expectedVegetable := "огурец"
	testName := "Test name"

	g := []products.Product{
		products.NewDefaultProduct("рис", products.Grain),
		products.NewDefaultProduct("гречка", products.Grain),
		products.NewDefaultProduct("овес", products.Grain),
	}
	p := []products.Product{
		products.NewDefaultProduct("яйцо", products.Protein),
		products.NewDefaultProduct("курица", products.Protein),
	}
	v := []products.Product{
		products.NewDefaultProduct("огурец", products.Vegetable),
		products.NewDefaultProduct("томат", products.Vegetable),
	}

	prod, err := GenerateDish(testName, g, p, v, rng)
	if err != nil {
		t.Fatalf("ошибка: %v", err)

	}
	if prod.Name != testName {
		t.Errorf("ожидали %q при сиде %d, получили %q", testName, seed, prod.Name)
	}
	if prod.Grain.Name != expectedGrain {
		t.Errorf("ожидали %q при сиде %d, получили %q", expectedGrain, seed, prod.Grain.Name)
	}
	if prod.Protein.Name != expectedProtein {
		t.Errorf("ожидали %q при сиде %d, получили %q", expectedProtein, seed, prod.Protein.Name)
	}
	if prod.Vegetable.Name != expectedVegetable {
		t.Errorf("ожидали %q при сиде %d, получили %q", expectedVegetable, seed, prod.Vegetable.Name)
	}

}
