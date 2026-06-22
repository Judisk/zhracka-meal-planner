package products

type ProductID int

type Product struct {
	ID       ProductID
	Name     string
	Category Category
}

type Category string

const (
	Grain     Category = "grain"
	Protein   Category = "protein"
	Vegetable Category = "vegetable"
)

var Grains = []Product{

	{ID: 1, Name: "рис", Category: Grain},
	{ID: 2, Name: "гречка", Category: Grain},
	{ID: 3, Name: "овес", Category: Grain},
}
var Proteins = []Product{
	{ID: 4, Name: "яйцо", Category: Protein},
	{ID: 5, Name: "курица", Category: Protein},
}
var Vegetables = []Product{
	{ID: 6, Name: "огурец", Category: Vegetable},
	{ID: 7, Name: "томат", Category: Vegetable},
}
