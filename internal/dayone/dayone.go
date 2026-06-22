package dayone

import (
	"fmt"
	f "foods/internal/foodgenerator"
	"foods/internal/products"
	"math/rand/v2"
)

type Day struct {
	Meals []f.Dish
}

func GenerateMeals(n int, grains, proteins, vegetables []products.Product, rng *rand.Rand) (Day, error) {
	arr := []f.Dish{}
	for i := 0; i < n; {

		dish, err := f.GenerateDish(namedMeals(i), grains, proteins, vegetables, rng)
		if err != nil {
			return Day{}, fmt.Errorf("generate meals: %w", err)
		}
		arr = append(arr, dish)
		i++
	}
	return Day{arr}, nil
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
