package dayone

import (
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

	g = append(g, products.NewProduct("banned grain", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned protein", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned vegetable", products.Vegetable, true, products.Neutral))

	g = append(g, products.NewProduct("banned grain1", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned protein1", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned vegetable1", products.Vegetable, true, products.Neutral))

	g = append(g, products.NewProduct("banned grain2", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned protein2", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned vegetable2", products.Vegetable, true, products.Neutral))

	g = append(g, products.NewProduct("banned grain3", products.Grain, true, products.Neutral))
	p = append(p, products.NewProduct("banned protein3", products.Protein, true, products.Neutral))
	v = append(v, products.NewProduct("banned vegetable3", products.Vegetable, true, products.Neutral))

	_, err := GenerateMeals(testN, g, p, v,
		products.BlockedProducts{
			ByName: map[string]bool{
				"banned grain": true, "banned protein": true, "banned vegetable": true,
				"banned grain1": true, "banned protein1": true, "banned vegetable1": true,
				"banned grain2": true, "banned protein2": true, "banned vegetable2": true,
				"banned grain3": true, "banned protein3": true, "banned vegetable3": true,
			},
		}, rng)

	if err == nil {
		t.Errorf("expected an error, got %v", err)
	}

}
