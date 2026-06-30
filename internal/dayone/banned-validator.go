package dayone

import (
	"fmt"
	"foods/internal/foodgenerator"
	"foods/internal/products"
)

func validator(dish foodgenerator.Dish, banned products.BlockedProducts) error {
	if _, ok := banned.ByID[dish.Grain.ID]; ok {
		return fmt.Errorf("found banned ID: %d", dish.Grain.ID)
	}
	if _, ok := banned.ByID[dish.Protein.ID]; ok {
		return fmt.Errorf("found banned ID: %d", dish.Protein.ID)
	}
	if _, ok := banned.ByID[dish.Vegetable.ID]; ok {
		return fmt.Errorf("found banned ID: %d", dish.Vegetable.ID)
	}
	if _, ok := banned.ByName[dish.Grain.Name]; ok {
		return fmt.Errorf("found banned name: %q", dish.Grain.Name)
	}
	if _, ok := banned.ByName[dish.Protein.Name]; ok {
		return fmt.Errorf("found banned name: %q", dish.Protein.Name)
	}
	if _, ok := banned.ByName[dish.Vegetable.Name]; ok {
		return fmt.Errorf("found banned name: %q", dish.Vegetable.Name)
	}

	return nil
}
