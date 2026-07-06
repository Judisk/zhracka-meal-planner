package service

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"foods/internal/storage"
)

func SeedDefaultProductsIfEmpty(db *sql.DB) error {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return fmt.Errorf("count products: %w", err)
	}

	if count == 0 {
		return InsertDefaultProducts(db)
	}

	return nil
}

func InsertDefaultProducts(db *sql.DB) error {

	defaults := []products.Product{
		// Grain / carbs
		products.NewDefaultProduct("rice", products.Grain),
		products.NewDefaultProduct("oats", products.Grain),
		products.NewDefaultProduct("buckwheat", products.Grain),
		products.NewDefaultProduct("bulgur", products.Grain),
		products.NewDefaultProduct("couscous", products.Grain),
		products.NewDefaultProduct("quinoa", products.Grain),
		products.NewDefaultProduct("millet", products.Grain),
		products.NewDefaultProduct("barley groats", products.Grain),
		products.NewDefaultProduct("semolina", products.Grain),
		products.NewDefaultProduct("cornmeal", products.Grain),
		products.NewDefaultProduct("spelt", products.Grain),
		products.NewDefaultProduct("amaranth", products.Grain),
		products.NewDefaultProduct("basmati rice", products.Grain),
		products.NewDefaultProduct("brown rice", products.Grain),
		products.NewDefaultProduct("wild rice", products.Grain),
		products.NewDefaultProduct("pasta", products.Grain),
		products.NewDefaultProduct("whole wheat pasta", products.Grain),
		products.NewDefaultProduct("noodles", products.Grain),
		products.NewDefaultProduct("buckwheat noodles", products.Grain),
		products.NewDefaultProduct("rice noodles", products.Grain),
		products.NewDefaultProduct("potatoes", products.Grain),
		products.NewDefaultProduct("sweet potato", products.Grain),
		products.NewDefaultProduct("lavash", products.Grain),
		products.NewDefaultProduct("whole grain bread", products.Grain),
		products.NewDefaultProduct("rye bread", products.Grain),
		products.NewDefaultProduct("pita", products.Grain),
		products.NewDefaultProduct("tortilla", products.Grain),
		products.NewDefaultProduct("glass noodles", products.Grain),
		products.NewDefaultProduct("pearl barley", products.Grain),

		// Protein
		products.NewDefaultProduct("egg", products.Protein),
		products.NewDefaultProduct("chicken", products.Protein),
		products.NewDefaultProduct("turkey", products.Protein),
		products.NewDefaultProduct("beef", products.Protein),
		products.NewDefaultProduct("pork", products.Protein),
		products.NewDefaultProduct("ground beef", products.Protein),
		products.NewDefaultProduct("ground chicken", products.Protein),
		products.NewDefaultProduct("tuna", products.Protein),
		products.NewDefaultProduct("salmon", products.Protein),
		products.NewDefaultProduct("cod", products.Protein),
		products.NewDefaultProduct("hake", products.Protein),
		products.NewDefaultProduct("mackerel", products.Protein),
		products.NewDefaultProduct("sardines", products.Protein),
		products.NewDefaultProduct("shrimp", products.Protein),
		products.NewDefaultProduct("squid", products.Protein),
		products.NewDefaultProduct("cottage cheese", products.Protein),
		products.NewDefaultProduct("greek yogurt", products.Protein),
		products.NewDefaultProduct("cheese", products.Protein),
		products.NewDefaultProduct("mozzarella", products.Protein),
		products.NewDefaultProduct("tofu", products.Protein),
		products.NewDefaultProduct("tempeh", products.Protein),
		products.NewDefaultProduct("beans", products.Protein),
		products.NewDefaultProduct("lentils", products.Protein),
		products.NewDefaultProduct("chickpeas", products.Protein),
		products.NewDefaultProduct("peas", products.Protein),
		products.NewDefaultProduct("soy", products.Protein),
		products.NewDefaultProduct("peanuts", products.Protein),
		products.NewDefaultProduct("almonds", products.Protein),
		products.NewDefaultProduct("chicken liver", products.Protein),
		products.NewDefaultProduct("ham", products.Protein),

		// Vegetable
		products.NewDefaultProduct("tomato", products.Vegetable),
		products.NewDefaultProduct("cucumber", products.Vegetable),
		products.NewDefaultProduct("carrot", products.Vegetable),
		products.NewDefaultProduct("onion", products.Vegetable),
		products.NewDefaultProduct("garlic", products.Vegetable),
		products.NewDefaultProduct("cabbage", products.Vegetable),
		products.NewDefaultProduct("napa cabbage", products.Vegetable),
		products.NewDefaultProduct("cauliflower", products.Vegetable),
		products.NewDefaultProduct("broccoli", products.Vegetable),
		products.NewDefaultProduct("zucchini", products.Vegetable),
		products.NewDefaultProduct("eggplant", products.Vegetable),
		products.NewDefaultProduct("bell pepper", products.Vegetable),
		products.NewDefaultProduct("chili pepper", products.Vegetable),
		products.NewDefaultProduct("beetroot", products.Vegetable),
		products.NewDefaultProduct("radish", products.Vegetable),
		products.NewDefaultProduct("daikon", products.Vegetable),
		products.NewDefaultProduct("pumpkin", products.Vegetable),
		products.NewDefaultProduct("celery", products.Vegetable),
		products.NewDefaultProduct("spinach", products.Vegetable),
		products.NewDefaultProduct("lettuce", products.Vegetable),
		products.NewDefaultProduct("arugula", products.Vegetable),
		products.NewDefaultProduct("green beans", products.Vegetable),
		products.NewDefaultProduct("green peas", products.Vegetable),
		products.NewDefaultProduct("corn", products.Vegetable),
		products.NewDefaultProduct("mushrooms", products.Vegetable),
		products.NewDefaultProduct("button mushrooms", products.Vegetable),
		products.NewDefaultProduct("asparagus", products.Vegetable),
		products.NewDefaultProduct("avocado", products.Vegetable),
		products.NewDefaultProduct("olives", products.Vegetable),
		products.NewDefaultProduct("herbs", products.Vegetable),
	}

	for _, p := range defaults {
		if err := storage.InsertProduct(db, p); err != nil {
			return fmt.Errorf("insert default product %q: %w", p.Name, err)
		}
	}

	return nil
}
