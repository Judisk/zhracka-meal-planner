package dayone

import (
	"fmt"
	"foods/internal/foodgenerator"
	"foods/internal/products"
)

func validator(dish foodgenerator.Dish, banned products.BlockedProducts) error {
	if ok := banned.ContainsID(dish.Grain.ID); ok {
		return fmt.Errorf("found banned ID: %d", dish.Grain.ID)
	}
	if ok := banned.ContainsID(dish.Protein.ID); ok {
		return fmt.Errorf("found banned ID: %d", dish.Protein.ID)
	}
	if ok := banned.ContainsID(dish.Vegetable.ID); ok {
		return fmt.Errorf("found banned ID: %d", dish.Vegetable.ID)
	}
	if ok := banned.ContainsName(dish.Grain.Name); ok {
		return fmt.Errorf("found banned name: %q", dish.Grain.Name)
	}
	if ok := banned.ContainsName(dish.Protein.Name); ok {
		return fmt.Errorf("found banned name: %q", dish.Protein.Name)
	}
	if ok := banned.ContainsName(dish.Vegetable.Name); ok {
		return fmt.Errorf("found banned name: %q", dish.Vegetable.Name)
	}

	return nil
}
