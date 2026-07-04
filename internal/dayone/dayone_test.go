package dayone

import (
	"errors"
	"foods/internal/products"
	"math/rand/v2"
	"testing"
)

func datasForTest() ([]string, []string, []string, []string, []products.Product, []products.Product, []products.Product) {
	arrayNamesMeals := []string{"Breakfast", "Lunch", "Dinner", "Snack 1", "Snack 2", "Snack 3"}

	expectedGrains := []string{"buckwheat", "oats", "rice", "tortilla", "pearl barley", "glass noodles"}
	expectedProteins := []string{"peanuts", "almonds", "chicken liver", "soy", "egg", "chicken"}
	expectedVegetables := []string{"cucumber", "tomato", "asparagus", "olives", "herbs", "avocado"}

	g := []products.Product{
		products.NewDefaultProduct("rice", products.Grain),
		products.NewDefaultProduct("buckwheat", products.Grain),
		products.NewDefaultProduct("oats", products.Grain),
		products.NewDefaultProduct("tortilla", products.Grain),
		products.NewDefaultProduct("glass noodles", products.Grain),
		products.NewDefaultProduct("pearl barley", products.Grain),
	}
	p := []products.Product{
		products.NewDefaultProduct("egg", products.Protein),
		products.NewDefaultProduct("chicken", products.Protein),
		products.NewDefaultProduct("soy", products.Protein),
		products.NewDefaultProduct("peanuts", products.Protein),
		products.NewDefaultProduct("almonds", products.Protein),
		products.NewDefaultProduct("chicken liver", products.Protein),
	}
	v := []products.Product{
		products.NewDefaultProduct("cucumber", products.Vegetable),
		products.NewDefaultProduct("tomato", products.Vegetable),
		products.NewDefaultProduct("asparagus", products.Vegetable),
		products.NewDefaultProduct("avocado", products.Vegetable),
		products.NewDefaultProduct("olives", products.Vegetable),
		products.NewDefaultProduct("herbs", products.Vegetable),
	}
	return arrayNamesMeals, expectedGrains, expectedProteins, expectedVegetables, g, p, v
}
func TestGenerateDay_FixedSeed_Success(t *testing.T) {
	var seed uint64 = 42
	var testN = 6
	rng := rand.New(rand.NewPCG(seed, 0))

	arrayNamesMeals, expectedGrains, expectedProteins, expectedVegetables, g, p, v := datasForTest()

	day, err := GenerateMeals(testN, g, p, v, products.BlockedProducts{}, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(day.Meals) != testN {
		t.Fatalf("wrong meal count: %v", err)
	}

	for i := range testN {
		if day.Meals[i].Name != arrayNamesMeals[i] {
			t.Errorf("expected %q with seed %d, got %q", arrayNamesMeals[i], seed, day.Meals[i].Name)
		}
		if day.Meals[i].Grain.Name != expectedGrains[i] {
			t.Errorf("%d expected %q with seed %d, got %q", i, expectedGrains[i], seed, day.Meals[i].Grain.Name)
		}
		if day.Meals[i].Protein.Name != expectedProteins[i] {
			t.Errorf("%d expected %q with seed %d, got %q", i, expectedProteins[i], seed, day.Meals[i].Protein.Name)
		}
		if day.Meals[i].Vegetable.Name != expectedVegetables[i] {
			t.Errorf("%d expected %q with seed %d, got %q", i, expectedVegetables[i], seed, day.Meals[i].Vegetable.Name)
		}
	}

}

func TestGenerateDay_FixedSeed_NExceedsProductCount(t *testing.T) {

	var seed uint64 = 42
	var testN = 7
	rng := rand.New(rand.NewPCG(seed, 0))

	_, _, _, _, g, p, v := datasForTest()

	_, err := GenerateMeals(testN, g, p, v, products.BlockedProducts{}, rng)
	if err == nil {
		t.Errorf("expected an error, got %v", err)
	}
}

func TestGenerateDay_EmptyLists(t *testing.T) {

	var seed uint64 = 42
	var testN = 1
	rng := rand.New(rand.NewPCG(seed, 0))
	g, p, v := []products.Product{}, []products.Product{}, []products.Product{}
	_, err := GenerateMeals(testN, g, p, v, products.BlockedProducts{}, rng)
	if err == nil {
		t.Errorf("expected an error, got %v", err)
	}

}

func TestGenerateDay_FixedSeed_GetBannedError(t *testing.T) {
	var seed uint64 = 42
	var testN = 10
	rng := rand.New(rand.NewPCG(seed, 0))

	g := []products.Product{}
	p := []products.Product{}
	v := []products.Product{}
	g = append(g, products.NewDefaultProduct("rice", products.Grain))
	p = append(p, products.NewDefaultProduct("egg", products.Protein))
	v = append(v, products.NewDefaultProduct("cucumber", products.Vegetable))

	g = append(g, products.NewProduct("banned Grain", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned Protein", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned Vegetable", products.Vegetable, true, products.Neutral))

	g = append(g, products.NewProduct("banned Grain1", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned Protein1", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned Vegetable1", products.Vegetable, true, products.Neutral))

	g = append(g, products.NewProduct("banned Grain2", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned Protein2", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned Vegetable2", products.Vegetable, true, products.Neutral))

	g = append(g, products.NewProduct("banned Grain3", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned Protein3", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned Vegetable3", products.Vegetable, true, products.Neutral))

	_, err := GenerateMeals(testN, g, p, v,
		products.BlockedProducts{
			ByName: map[string]bool{
				"banned Grain": true, "banned Protein": true, "banned Vegetable": true,
				"banned Grain1": true, "banned Protein1": true, "banned Vegetable1": true,
				"banned Grain2": true, "banned Protein2": true, "banned Vegetable2": true,
				"banned Grain3": true, "banned Protein3": true, "banned Vegetable3": true,
			},
		}, rng)

	if err == nil {
		t.Errorf("expected an error, got %v", err)
	}
	if !errors.Is(err, ErrTooManyErrors) {
		t.Fatalf("expected an error %v, got %v", ErrTooManyErrors, err)
	}

}
