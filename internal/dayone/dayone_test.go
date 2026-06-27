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

	expectedGrains := []string{"рис", "рис", "рис", "рис"}
	expectedProteins := []string{"курица", "курица", "курица", "курица"}
	expectedVegetables := []string{"огурец", "огурец", "огурец", "огурец"}

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

	day, err := GenerateMeals(testN, g, p, v, rng)
	if err != nil {
		t.Fatalf("ошибка: %v", err)
	}
	if len(day.Meals) != testN {
		t.Fatalf("ошибка длины массива: %v", err)
	}

	for i := range testN {
		if day.Meals[i].Name != arrayNamesMeals[i] {
			t.Errorf("ожидали %q при сиде %d, получили %q", arrayNamesMeals[i], seed, day.Meals[i].Name)
		}
		if day.Meals[i].Grain.Name != expectedGrains[i] {
			t.Errorf("%d ожидали %q при сиде %d, получили %q", i, expectedGrains[i], seed, day.Meals[i].Grain.Name)
		}
		if day.Meals[i].Protein.Name != expectedProteins[i] {
			t.Errorf("%d ожидали %q при сиде %d, получили %q", i, expectedProteins[i], seed, day.Meals[i].Protein.Name)
		}
		if day.Meals[i].Vegetable.Name != expectedVegetables[i] {
			t.Errorf("%d ожидали %q при сиде %d, получили %q", i, expectedVegetables[i], seed, day.Meals[i].Vegetable.Name)
		}
	}

}
