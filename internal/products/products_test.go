package products

import (
	"testing"
)

func TestNormalizeProductName_Table(t *testing.T) {
	tests := []struct {
		name         string
		testText     string
		expectedText string
	}{
		{
			name:         "happy path",
			testText:     "Rice",
			expectedText: "rice",
		},
		{
			name:         "already normalize",
			testText:     "rice",
			expectedText: "rice",
		},
		{
			name:         "empty string",
			testText:     "",
			expectedText: "",
		},
		{
			name:         "space only",
			testText:     "              ",
			expectedText: "",
		},
		{
			name:         "trim and lowercase",
			testText:     "  RICE  ",
			expectedText: "rice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeProductName(tt.testText)
			if result != tt.expectedText {
				t.Errorf("expected an %q got %q", tt.expectedText, result)
			}
		})
	}
}

func TestNewBlockedProducts_Table(t *testing.T) {
	type expectedId struct {
		id     ProductID
		wanted bool
	}
	type expectedName struct {
		name   string
		wanted bool
	}

	tests := []struct {
		name     string
		mapIDs   map[ProductID]bool
		rawNames []string

		expectedIDs   []expectedId
		expectedNames []expectedName
	}{
		{
			name:     "happy",
			mapIDs:   map[ProductID]bool{1: true, 2: true},
			rawNames: []string{"rice", "chicken"},
			expectedIDs: []expectedId{
				{id: 1, wanted: true}, {id: 2, wanted: true},
				{id: 4, wanted: false},
			},
			expectedNames: []expectedName{{name: "rice", wanted: true},
				{name: "chicken", wanted: true},
				{name: "potato", wanted: false}},
		},
		{
			name:     "big names",
			mapIDs:   map[ProductID]bool{1: true, 2: true},
			rawNames: []string{"RICE", "CHICKEN"},
			expectedIDs: []expectedId{
				{id: 1, wanted: true}, {id: 2, wanted: true},
				{id: 4, wanted: false},
			},
			expectedNames: []expectedName{{name: "rice", wanted: true},
				{name: "chicken", wanted: true},
				{name: "potato", wanted: false}},
		},
		{
			name:     "have id empty names",
			mapIDs:   map[ProductID]bool{1: true, 2: true},
			rawNames: []string{},
			expectedIDs: []expectedId{
				{id: 1, wanted: true}, {id: 2, wanted: true},
				{id: 4, wanted: false},
			},
			expectedNames: []expectedName{{name: "potato", wanted: false}},
		},
		{
			name:        "nil id map returns false, not panic",
			mapIDs:      nil,
			rawNames:    []string{"rice", "chicken"},
			expectedIDs: []expectedId{{id: 4, wanted: false}},
			expectedNames: []expectedName{{name: "rice", wanted: true},
				{name: "chicken", wanted: true},
				{name: "potato", wanted: false}},
		},
		{
			name:          "empty map and empty names",
			mapIDs:        map[ProductID]bool{},
			rawNames:      []string{},
			expectedIDs:   []expectedId{{id: 4, wanted: false}},
			expectedNames: []expectedName{{name: "potato", wanted: false}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blockedProducts := NewBlockedProducts(tt.mapIDs, tt.rawNames)
			for _, elem := range tt.expectedIDs {
				got := blockedProducts.ContainsID(elem.id)
				if got != elem.wanted {
					t.Errorf("ContainsID(%d) = %t, want %t", elem.id, got, elem.wanted)
				}
			}
			for _, elem := range tt.expectedNames {
				got := blockedProducts.ContainsName(elem.name)
				if got != elem.wanted {
					t.Errorf("ContainsName(%q) = %t, want %t", elem.name, got, elem.wanted)
				}
			}

		})

	}
}

func TestNewProducts_Table(t *testing.T) {

	type testedProd struct {
		prod   Product
		wanted Product
	}
	tests := []struct {
		name           string
		testedProducts []testedProd
		fn             func(testedProd) Product
	}{
		{
			name: "default prod",
			testedProducts: []testedProd{
				{prod: Product{
					Name: "rice", Category: Grain},
					wanted: Product{
						Name:           "rice",
						Category:       Grain,
						Banned:         false,
						Preference:     Neutral,
						SelectionScore: 1.0,
					}},
				{prod: Product{
					Name: "egg", Category: Protein},
					wanted: Product{
						Name:           "egg",
						Category:       Protein,
						Banned:         false,
						Preference:     Neutral,
						SelectionScore: 1.0}},
				{prod: Product{
					Name: "tomato", Category: Vegetable},
					wanted: Product{
						Name:           "tomato",
						Category:       Vegetable,
						Banned:         false,
						Preference:     Neutral,
						SelectionScore: 1.0,
					}},
			},
			fn: func(testedProducts testedProd) Product {
				p := testedProducts.prod
				return NewDefaultProduct(p.Name, p.Category)
			},
		},
		{
			name: "non-default prod",
			testedProducts: []testedProd{
				{
					prod: Product{
						Name:       "rice",
						Category:   Grain,
						Banned:     false,
						Preference: Neutral,
					},
					wanted: Product{
						Name:           "rice",
						Category:       Grain,
						Banned:         false,
						Preference:     Neutral,
						SelectionScore: 1.0,
					}},
				{
					prod: Product{
						Name:       "egg",
						Category:   Protein,
						Banned:     false,
						Preference: Liked,
					},
					wanted: Product{
						Name:           "egg",
						Category:       Protein,
						Banned:         false,
						Preference:     Liked,
						SelectionScore: 1.0}},
				{
					prod: Product{
						Name:       "tomato",
						Category:   Vegetable,
						Banned:     true,
						Preference: Disliked,
					},
					wanted: Product{
						Name:           "tomato",
						Category:       Vegetable,
						Banned:         true,
						Preference:     Disliked,
						SelectionScore: 1.0,
					}},
			},
			fn: func(testedProducts testedProd) Product {
				p := testedProducts.prod
				return NewProduct(p.Name, p.Category, p.Banned, p.Preference)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, elem := range tt.testedProducts {
				gotProduct := tt.fn(elem)
				if gotProduct.Name != elem.wanted.Name {
					t.Errorf("unexpected name %q wanted %q", gotProduct.Name, elem.wanted.Name)
				}
				if gotProduct.Category != elem.wanted.Category {
					t.Errorf("unexpected category %q wanted %q", gotProduct.Category, elem.wanted.Category)
				}
				if gotProduct.Banned != elem.wanted.Banned {
					t.Errorf("unexpected banned %t wanted %t", gotProduct.Banned, elem.wanted.Banned)
				}
				if gotProduct.Preference != elem.wanted.Preference {
					t.Errorf("unexpected preference %f wanted %f", float64(gotProduct.Preference), float64(elem.wanted.Preference))
				}
				if gotProduct.SelectionScore != elem.wanted.SelectionScore {
					t.Errorf("unexpected SelectionScore %f wanted %f", gotProduct.SelectionScore, elem.wanted.SelectionScore)
				}
			}

		})
	}
}
