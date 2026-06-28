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

func TestGenerateDay_FixedSeed_NExceedsProductCount(t *testing.T) {

	var seed uint64 = 42
	var testN = 7
	rng := rand.New(rand.NewPCG(seed, 0))

	_, _, _, _, g, p, v := datasForTest()

	_, err := GenerateMeals(testN, g, p, v, rng)
	if err == nil {
		t.Errorf("ожидалась ошибка получили %v", err)
	}

}

func TestGenerateDay_EmptyLists(t *testing.T) {

	var seed uint64 = 42
	var testN = 1
	rng := rand.New(rand.NewPCG(seed, 0))
	g, p, v := []products.Product{}, []products.Product{}, []products.Product{}
	_, err := GenerateMeals(testN, g, p, v, rng)
	if err == nil {
		t.Errorf("ожидалась ошибка получили %v", err)
	}
}
