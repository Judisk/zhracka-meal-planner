package products

type ProductID int

type Product struct {
	ID             ProductID
	Name           string
	Category       Category
	Banned         bool
	Preference     PreferenceStatus
	SelectionScore float64
}

type Category string

const (
	Grain     Category = "grain"
	Protein   Category = "protein"
	Vegetable Category = "vegetable"
)

type PreferenceStatus float64

const (
	Liked    PreferenceStatus = 1.5
	Neutral  PreferenceStatus = 1
	Disliked PreferenceStatus = 0.5
)

type BlockedProducts struct {
	ByID   map[ProductID]bool
	ByName map[string]bool
}

func NewProduct(name string, category Category) Product {
	return Product{
		Name:           name,
		Category:       category,
		Banned:         false,
		Preference:     Neutral,
		SelectionScore: 1.0,
	}
}
