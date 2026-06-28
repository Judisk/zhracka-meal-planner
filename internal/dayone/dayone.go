package dayone

import (
	"fmt"
	f "foods/internal/foodgenerator"
	"foods/internal/products"
	"math/rand/v2"
)

const maxErrorCount int = 10

type Day struct {
	Meals []f.Dish
}

func GenerateMeals(n int, grains, proteins, vegetables []products.Product, banned products.BlockedProducts, rng *rand.Rand) (Day, error) {
	prods := []f.Dish{}
	errCount := 0
	for i := 0; i < n; {

		dish, err := f.GenerateDish(namedMeals(i), grains, proteins, vegetables, rng)
		if err != nil {
			return Day{}, fmt.Errorf("generate meals: %w", err)
		}
		if err := validator(dish, banned); err != nil {
			if errCount == maxErrorCount {
				return Day{}, fmt.Errorf("generate meals: too many errors: %w", err)
			}
			errCount++
			continue
		}
		grains, proteins, vegetables = filterUsedProducts(grains, proteins, vegetables, dish)
		prods = append(prods, dish)
		i++
	}
	return Day{prods}, nil
}

func namedMeals(n int) string {
	switch n {
	case 0:
		return "Завтрак"
	case 1:
		return "Обед"
	case 2:
		return "Ужин"
	default:
		return fmt.Sprintf("Перекус %d", n-2)
	}

}

func filterUsedProducts(grains, proteins, vegetables []products.Product, dish f.Dish) ([]products.Product, []products.Product, []products.Product) {

	arrOne := filterOut(grains, dish.Grain)
	arrTwo := filterOut(proteins, dish.Protein)
	arrThree := filterOut(vegetables, dish.Vegetable)
	return arrOne, arrTwo, arrThree
}

func filterOut(array []products.Product, pr products.Product) []products.Product {
	result := make([]products.Product, 0, len(array)-1)
	for _, p := range array {
		if p != pr {
			result = append(result, p)
		}
	}
	return result
}
