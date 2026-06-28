package dayone

import (
	"foods/internal/products"
	"math/rand/v2"
	"testing"
)

func datasForTest() ([]string, []string, []string, []string, []products.Product, []products.Product, []products.Product) {
	arrayNamesMeals := []string{"Завтрак", "Обед", "Ужин", "Перекус 1", "Перекус 2", "Перекус 3"}

	expectedGrains := []string{"гречка", "овес", "рис", "тортилья", "перловая крупа", "фунчоза"}
	expectedProteins := []string{"арахис", "миндаль", "куриная печень", "соя", "яйцо", "курица"}
	expectedVegetables := []string{"огурец", "томат", "спаржа", "оливки", "зелень", "авокадо"}

	g := []products.Product{
		products.NewDefaultProduct("рис", products.Grain),
		products.NewDefaultProduct("гречка", products.Grain),
		products.NewDefaultProduct("овес", products.Grain),
		products.NewDefaultProduct("тортилья", products.Grain),
		products.NewDefaultProduct("фунчоза", products.Grain),
		products.NewDefaultProduct("перловая крупа", products.Grain),
	}
	p := []products.Product{
		products.NewDefaultProduct("яйцо", products.Protein),
		products.NewDefaultProduct("курица", products.Protein),
		products.NewDefaultProduct("соя", products.Protein),
		products.NewDefaultProduct("арахис", products.Protein),
		products.NewDefaultProduct("миндаль", products.Protein),
		products.NewDefaultProduct("куриная печень", products.Protein),
	}
	v := []products.Product{
		products.NewDefaultProduct("огурец", products.Vegetable),
		products.NewDefaultProduct("томат", products.Vegetable),
		products.NewDefaultProduct("спаржа", products.Vegetable),
		products.NewDefaultProduct("авокадо", products.Vegetable),
		products.NewDefaultProduct("оливки", products.Vegetable),
		products.NewDefaultProduct("зелень", products.Vegetable),
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

func TestGenerateDay_FixedSeed_NExceedsProductCount(t *testing.T) {

	var seed uint64 = 42
	var testN = 7
	rng := rand.New(rand.NewPCG(seed, 0))

	_, _, _, _, g, p, v := datasForTest()

	_, err := GenerateMeals(testN, g, p, v, products.BlockedProducts{}, rng)
	if err == nil {
		t.Errorf("ожидалась ошибка получили %v", err)
	}

}

func TestGenerateDay_EmptyLists(t *testing.T) {

	var seed uint64 = 42
	var testN = 1
	rng := rand.New(rand.NewPCG(seed, 0))
	g, p, v := []products.Product{}, []products.Product{}, []products.Product{}
	_, err := GenerateMeals(testN, g, p, v, products.BlockedProducts{}, rng)
	if err == nil {
		t.Errorf("ожидалась ошибка получили %v", err)
	}
}

func TestGenerateDay_FixedSeed_GetBannedError(t *testing.T) {
	var seed uint64 = 42
	var testN = 10
	rng := rand.New(rand.NewPCG(seed, 0))

	g := []products.Product{}
	p := []products.Product{}
	v := []products.Product{}
	g = append(g, products.NewDefaultProduct("рис", products.Grain))
	p = append(p, products.NewDefaultProduct("яйцо", products.Protein))
	v = append(v, products.NewDefaultProduct("огурец", products.Vegetable))

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
		t.Errorf("ожидалась ошибка получили %v", err)
	}

}
