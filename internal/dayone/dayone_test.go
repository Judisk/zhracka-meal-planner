package dayone

import (
	"foods/internal/products"
	"math/rand/v2"
	"testing"
)

func TestGenerateDay_FixedSeed(t *testing.T) {
	var seed uint64 = 42
	var testN = 4
	rng := rand.New(rand.NewPCG(seed, 0))
	arrayNamesMeals := []string{"Завтрак", "Обед", "Ужин", "Перекус 1"}

	expectedGrains := []string{"овес", "рис", "овес", "овес"}
	expectedProteins := []string{"яйцо", "курица", "яйцо", "яйцо"}
	expectedVegetables := []string{"томат", "томат", "томат", "огурец"}

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

	day, err := GenerateMeals(testN, g, p, v, rng)
	if err != nil {
		t.Fatalf("ошибка: %v", err)
	}
	if len(day.Meals) != testN {
		t.Fatalf("ошибка длины массива: %v", err)
	}

	for i := range testN {
		if day.Meals[i].Name != arrayNamesMeals[i] {
			t.Errorf("ожидали '%s' при сиде %d, получили %q", arrayNamesMeals[i], seed, day.Meals[i].Name)
		}
		if day.Meals[i].Grain.Name != expectedGrains[i] {
			t.Errorf("%d ожидали '%s' при сиде %d, получили %q", i, expectedGrains[i], seed, day.Meals[i].Grain.Name)
		}
		if day.Meals[i].Protein.Name != expectedProteins[i] {
			t.Errorf("%d ожидали '%s' при сиде %d, получили %q", i, expectedProteins[i], seed, day.Meals[i].Protein.Name)
		}
		if day.Meals[i].Vegetable.Name != expectedVegetables[i] {
			t.Errorf("%d ожидали '%s' при сиде %d, получили %q", i, expectedVegetables[i], seed, day.Meals[i].Vegetable.Name)
		}
	}

}
