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
		return insertDefaultProducts(db)
	}

	return nil
}

func insertDefaultProducts(db *sql.DB) error {

	defaults := []products.Product{
		// Grain / carbs
		products.NewDefaultProduct("рис", products.Grain),
		products.NewDefaultProduct("овес", products.Grain),
		products.NewDefaultProduct("гречка", products.Grain),
		products.NewDefaultProduct("булгур", products.Grain),
		products.NewDefaultProduct("кускус", products.Grain),
		products.NewDefaultProduct("киноа", products.Grain),
		products.NewDefaultProduct("пшено", products.Grain),
		products.NewDefaultProduct("ячневая крупа", products.Grain),
		products.NewDefaultProduct("манка", products.Grain),
		products.NewDefaultProduct("кукурузная крупа", products.Grain),
		products.NewDefaultProduct("полба", products.Grain),
		products.NewDefaultProduct("амарант", products.Grain),
		products.NewDefaultProduct("рис басмати", products.Grain),
		products.NewDefaultProduct("бурый рис", products.Grain),
		products.NewDefaultProduct("дикий рис", products.Grain),
		products.NewDefaultProduct("макароны", products.Grain),
		products.NewDefaultProduct("цельнозерновые макароны", products.Grain),
		products.NewDefaultProduct("лапша", products.Grain),
		products.NewDefaultProduct("гречневая лапша", products.Grain),
		products.NewDefaultProduct("рисовая лапша", products.Grain),
		products.NewDefaultProduct("картофель", products.Grain),
		products.NewDefaultProduct("батат", products.Grain),
		products.NewDefaultProduct("лаваш", products.Grain),
		products.NewDefaultProduct("цельнозерновой хлеб", products.Grain),
		products.NewDefaultProduct("ржаной хлеб", products.Grain),
		products.NewDefaultProduct("пита", products.Grain),
		products.NewDefaultProduct("тортилья", products.Grain),
		products.NewDefaultProduct("фунчоза", products.Grain),
		products.NewDefaultProduct("перловая крупа", products.Grain),

		// Protein
		products.NewDefaultProduct("яйцо", products.Protein),
		products.NewDefaultProduct("курица", products.Protein),
		products.NewDefaultProduct("индейка", products.Protein),
		products.NewDefaultProduct("говядина", products.Protein),
		products.NewDefaultProduct("свинина", products.Protein),
		products.NewDefaultProduct("фарш говяжий", products.Protein),
		products.NewDefaultProduct("фарш куриный", products.Protein),
		products.NewDefaultProduct("тунец", products.Protein),
		products.NewDefaultProduct("лосось", products.Protein),
		products.NewDefaultProduct("треска", products.Protein),
		products.NewDefaultProduct("хек", products.Protein),
		products.NewDefaultProduct("скумбрия", products.Protein),
		products.NewDefaultProduct("сардины", products.Protein),
		products.NewDefaultProduct("креветки", products.Protein),
		products.NewDefaultProduct("кальмар", products.Protein),
		products.NewDefaultProduct("творог", products.Protein),
		products.NewDefaultProduct("греческий йогурт", products.Protein),
		products.NewDefaultProduct("сыр", products.Protein),
		products.NewDefaultProduct("моцарелла", products.Protein),
		products.NewDefaultProduct("тофу", products.Protein),
		products.NewDefaultProduct("темпе", products.Protein),
		products.NewDefaultProduct("фасоль", products.Protein),
		products.NewDefaultProduct("чечевица", products.Protein),
		products.NewDefaultProduct("нут", products.Protein),
		products.NewDefaultProduct("горох", products.Protein),
		products.NewDefaultProduct("соя", products.Protein),
		products.NewDefaultProduct("арахис", products.Protein),
		products.NewDefaultProduct("миндаль", products.Protein),
		products.NewDefaultProduct("куриная печень", products.Protein),
		products.NewDefaultProduct("ветчина", products.Protein),

		// Vegetable
		products.NewDefaultProduct("томат", products.Vegetable),
		products.NewDefaultProduct("огурец", products.Vegetable),
		products.NewDefaultProduct("морковь", products.Vegetable),
		products.NewDefaultProduct("лук", products.Vegetable),
		products.NewDefaultProduct("чеснок", products.Vegetable),
		products.NewDefaultProduct("капуста", products.Vegetable),
		products.NewDefaultProduct("пекинская капуста", products.Vegetable),
		products.NewDefaultProduct("цветная капуста", products.Vegetable),
		products.NewDefaultProduct("брокколи", products.Vegetable),
		products.NewDefaultProduct("кабачок", products.Vegetable),
		products.NewDefaultProduct("баклажан", products.Vegetable),
		products.NewDefaultProduct("болгарский перец", products.Vegetable),
		products.NewDefaultProduct("острый перец", products.Vegetable),
		products.NewDefaultProduct("свекла", products.Vegetable),
		products.NewDefaultProduct("редис", products.Vegetable),
		products.NewDefaultProduct("дайкон", products.Vegetable),
		products.NewDefaultProduct("тыква", products.Vegetable),
		products.NewDefaultProduct("сельдерей", products.Vegetable),
		products.NewDefaultProduct("шпинат", products.Vegetable),
		products.NewDefaultProduct("салат", products.Vegetable),
		products.NewDefaultProduct("руккола", products.Vegetable),
		products.NewDefaultProduct("зеленая фасоль", products.Vegetable),
		products.NewDefaultProduct("зеленый горошек", products.Vegetable),
		products.NewDefaultProduct("кукуруза", products.Vegetable),
		products.NewDefaultProduct("грибы", products.Vegetable),
		products.NewDefaultProduct("шампиньоны", products.Vegetable),
		products.NewDefaultProduct("спаржа", products.Vegetable),
		products.NewDefaultProduct("авокадо", products.Vegetable),
		products.NewDefaultProduct("оливки", products.Vegetable),
		products.NewDefaultProduct("зелень", products.Vegetable),
	}

	for _, p := range defaults {
		if err := storage.InsertProduct(db, p); err != nil {
			return fmt.Errorf("insert default product %q: %w", p.Name, err)
		}
	}

	return nil
}
