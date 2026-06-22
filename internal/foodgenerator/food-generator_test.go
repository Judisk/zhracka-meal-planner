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

	array := []products.Product{
		{Name: "Рис", Category: products.Grain},
		{Name: "Гречка", Category: products.Grain},
		{Name: "Овес", Category: products.Grain},
	}
	prod, err := getProd(array, rng)
	if err != nil {
		t.Fatalf("ошибка: %v", err)

	}
	if prod.Name != "Овес" {
		t.Errorf("ожидали 'Овес' при сиде 42, получили %q", prod.Name)
	}
}

func TestGenerateDish_FixedSeed(t *testing.T) {
	var seed uint64 = 42
	rng := rand.New(rand.NewPCG(seed, 0))

	expectedGrain := "овес"
	expectedProtein := "яйцо"
	expectedVegetable := "томат"
	testName := "Test name"

	g := []products.Product{

		{Name: "рис", Category: products.Grain},
		{Name: "гречка", Category: products.Grain},
		{Name: "овес", Category: products.Grain},
	}
	p := []products.Product{
		{Name: "яйцо", Category: products.Protein},
		{Name: "курица", Category: products.Protein},
	}
	v := []products.Product{
		{Name: "огурец", Category: products.Vegetable},
		{Name: "томат", Category: products.Vegetable},
	}

	prod, err := GenerateDish(testName, g, p, v, rng)
	if err != nil {
		t.Fatalf("ошибка: %v", err)

	}
	if prod.Name != testName {
		t.Errorf("ожидали '%s' при сиде %d, получили %q", testName, seed, prod.Name)
	}
	if prod.Grain.Name != expectedGrain {
		t.Errorf("ожидали '%s' при сиде %d, получили %q", expectedGrain, seed, prod.Grain.Name)
	}
	if prod.Protein.Name != expectedProtein {
		t.Errorf("ожидали '%s' при сиде %d, получили %q", expectedProtein, seed, prod.Protein.Name)
	}
	if prod.Vegetable.Name != expectedVegetable {
		t.Errorf("ожидали '%s' при сиде %d, получили %q", expectedVegetable, seed, prod.Vegetable.Name)
	}

}
