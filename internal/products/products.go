package products

import (
	"strings"
)

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
	Grain     Category = "Grain"
	Protein   Category = "Protein"
	Vegetable Category = "Vegetable"
)

type PreferenceStatus float64

const (
	Liked    PreferenceStatus = 1.5
	Neutral  PreferenceStatus = 1
	Disliked PreferenceStatus = 0.5
)

type BlockedProducts struct {
	byID   map[ProductID]bool
	byName map[string]bool
}

func (b BlockedProducts) ContainsID(id ProductID) bool {
	return b.byID[id]
}
func (b BlockedProducts) ContainsName(name string) bool {
	return b.byName[NormalizeProductName(name)]
}
func NewBlockedProducts(bannedByID map[ProductID]bool, rawNames []string) BlockedProducts {
	byName := make(map[string]bool)
	for _, name := range rawNames {
		byName[NormalizeProductName(name)] = true
	}
	return BlockedProducts{byID: bannedByID, byName: byName}
}

func NewDefaultProduct(name string, category Category) Product {
	return NewProduct(name, category, false, Neutral)
}

func NewProduct(name string, category Category, banned bool, preference PreferenceStatus) Product {
	return Product{
		Name:           name,
		Category:       category,
		Banned:         banned,
		Preference:     preference,
		SelectionScore: 1.0,
	}
}

func NormalizeProductName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
