package service

import (
	"database/sql"
	"errors"
	"foods/internal/products"
	"foods/internal/storage"
	"math/rand/v2"
	"testing"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := storage.NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return db
}

func setupDBAndRand(t *testing.T) (*sql.DB, *rand.Rand) {
	t.Helper()
	return setupDB(t), rand.New(rand.NewPCG(42, 0))
}
func insertDB(t *testing.T, db *sql.DB, p products.Product) {
	t.Helper()
	err := storage.InsertProduct(db, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
func productFounder(t *testing.T, db *sql.DB, id products.ProductID) ProdsForGui {
	t.Helper()
	p, err := storage.SelectProductByID(db, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return ProdsForGui{p}
}

func TestDayControl_NOutOfCount(t *testing.T) {
	tests := []struct {
		name        string
		n           int
		expectedErr error
	}{{
		name:        "n exceeds maximum",
		n:           7,
		expectedErr: ErrTooManyMeals,
	}, {
		name:        "n is below minimum",
		n:           0,
		expectedErr: ErrTooFewMeals,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, rng := setupDBAndRand(t)
			defer db.Close()
			_, err := GenerateAndControlDay(db, tt.n, rng)
			if err == nil {
				t.Fatalf("expected an error, got %v", err)
			}
			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("expected an error %v, got %v", tt.expectedErr, err)
			}
		})

	}
}
func TestDayControl_GenerateAndControlDayWithProds(t *testing.T) {
	tests := []struct {
		name            string
		expectedAnError bool
		n               int
	}{{
		name:            "GenerateAndControlDay Success",
		expectedAnError: false,
		n:               3,
	}, {
		name:            "GenerateAndControlDay N More Than Prods",
		expectedAnError: true,
		n:               4,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, rng := setupDBAndRand(t)
			defer db.Close()
			N := tt.n
			usedProds := make(map[string]bool)
			prods := []products.Product{
				products.NewDefaultProduct("rice", products.Grain),
				products.NewDefaultProduct("buckwheat", products.Grain),
				products.NewDefaultProduct("oats", products.Grain),

				products.NewDefaultProduct("egg", products.Protein),
				products.NewDefaultProduct("chicken", products.Protein),
				products.NewDefaultProduct("ham", products.Protein),

				products.NewDefaultProduct("herbs", products.Vegetable),
				products.NewDefaultProduct("cucumber", products.Vegetable),
				products.NewDefaultProduct("tomato", products.Vegetable),
			}

			for _, p := range prods {
				if err := storage.InsertProduct(db, p); err != nil {
					t.Fatalf("unexpected error %v", err)
				}
			}

			day, err := GenerateAndControlDay(db, N, rng)

			if tt.expectedAnError {
				if err == nil {
					t.Fatalf("expected an error got %v", err)
				}

			} else {
				if err != nil {
					t.Fatalf("unexpected error %v", err)
				}

				for _, elem := range day.Meals {
					if _, ok := usedProds[elem.Grain.Name]; ok {
						t.Errorf("unexpected repeat %q", elem.Grain.Name)
					}

					if _, ok := usedProds[elem.Vegetable.Name]; ok {
						t.Errorf("unexpected repeat %q", elem.Vegetable.Name)
					}

					if _, ok := usedProds[elem.Protein.Name]; ok {
						t.Errorf("unexpected repeat %q", elem.Protein.Name)
					}
					usedProds[elem.Grain.Name] = true
					usedProds[elem.Vegetable.Name] = true
					usedProds[elem.Protein.Name] = true
				}
			}
		})
	}
}

func TestDayControl_GenerateAndControlDayEmptyDB(t *testing.T) {
	db, rng := setupDBAndRand(t)
	defer db.Close()
	N := 3

	_, err := GenerateAndControlDay(db, N, rng)
	if err == nil {
		t.Fatalf("expected an error got %v", err)
	}

}

func TestService_ClosedDB_Table(t *testing.T) {

	type setup struct {
		n         int
		prod      ProdsForGui
		wasBanned bool
		rng       *rand.Rand
	}
	tests := []struct {
		name string
		st   setup
		fn   func(*sql.DB, setup) error
	}{
		{
			name: "add",
			st: setup{
				prod: ProdsForGui{products.NewDefaultProduct("rice", products.Grain)},
			},
			fn: func(db *sql.DB, p setup) error {
				return p.prod.Add(db)

			},
		},
		{
			name: "edit",
			st: setup{
				prod:      ProdsForGui{products.NewDefaultProduct("rice", products.Grain)},
				wasBanned: false,
			},
			fn: func(db *sql.DB, p setup) error {
				return p.prod.Edit(db, p.wasBanned)

			},
		},
		{
			name: "delete",
			st: setup{
				prod: ProdsForGui{products.NewDefaultProduct("rice", products.Grain)},
			},
			fn: func(db *sql.DB, p setup) error {
				return p.prod.Delete(db)

			},
		},
		{
			name: "Generate And Control",
			st: setup{
				n:   3,
				rng: rand.New(rand.NewPCG(42, 0)),
			},
			fn: func(db *sql.DB, p setup) error {
				_, err := GenerateAndControlDay(db, p.n, p.rng)
				return err

			},
		},
		{
			name: "get list",
			fn: func(db *sql.DB, s setup) error {
				_, err := GetList(db)
				return err
			},
		},
		{
			name: "get list filtered",
			fn: func(db *sql.DB, s setup) error {
				_, err := GetListFiltered(db, nil, nil, nil)
				return err
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db := setupDB(t)

			db.Close()

			err := tt.fn(db, tt.st)
			if err == nil {
				t.Errorf("expected an error got %v", err)
			}
		})
	}
}

func TestServiceProds_Edit_Table(t *testing.T) {

	type testedProd struct {
		oldProd      ProdsForGui
		newProd      ProdsForGui
		expectedProd ProdsForGui
	}

	tests := []struct {
		name          string
		expectedError bool
		testProd      testedProd
	}{
		{
			name:          "change preference",
			expectedError: false,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "rice",
						Category:   products.Grain,
						Banned:     false,
						Preference: products.Liked,
					}},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Liked,
						SelectionScore: 14.0,
					}},
			}},
		{
			name:          "change name",
			expectedError: false,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "rices",
						Category:   products.Grain,
						Banned:     false,
						Preference: products.Neutral,
					}},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "rices",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					}},
			}},
		{
			name:          "change category",
			expectedError: false,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "rice",
						Category:   products.Vegetable,
						Banned:     false,
						Preference: products.Neutral,
					}},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Vegetable,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					}},
			}},
		{
			name:          "change banned",
			expectedError: false,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "rice",
						Category:   products.Grain,
						Banned:     true,
						Preference: products.Neutral,
					}},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         true,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					}},
			}},
		{
			name:          "unban resets score to 1.0",
			expectedError: false,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         true,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "rice",
						Category:   products.Grain,
						Banned:     false,
						Preference: products.Neutral,
					}},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 1.0,
					}},
			}},
		{
			name:          "change name,category,banned and preference",
			expectedError: false,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "soy",
						Category:   products.Protein,
						Banned:     true,
						Preference: products.Disliked,
					}},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "soy",
						Category:       products.Protein,
						Banned:         true,
						Preference:     products.Disliked,
						SelectionScore: 14.0,
					}},
			}},
		{
			name:          "change all",
			expectedError: false,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         true,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "soy",
						Category:   products.Protein,
						Banned:     false,
						Preference: products.Disliked,
					}},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "soy",
						Category:       products.Protein,
						Banned:         false,
						Preference:     products.Disliked,
						SelectionScore: 1.0,
					}},
			}},
		{
			name:          "error empty name",
			expectedError: true,
			testProd: testedProd{
				oldProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         true,
						Preference:     products.Neutral,
						SelectionScore: 14.0,
					},
				},
				newProd: ProdsForGui{
					products.Product{
						Name:       "",
						Category:   products.Protein,
						Banned:     false,
						Preference: products.Disliked,
					}},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var foundId products.ProductID = 1
			old := tt.testProd.oldProd
			expect := tt.testProd.expectedProd
			db := setupDB(t)
			insertDB(t, db, old.Prod)

			tt.testProd.newProd.Prod.ID = foundId
			err := tt.testProd.newProd.Edit(db, old.Prod.Banned)
			if (err != nil) != tt.expectedError {
				t.Fatalf("for error %v expected %t, got %t", err, tt.expectedError, (err != nil))
			}
			if !tt.expectedError {

				new := productFounder(t, db, foundId)

				if expect.Prod.Name != new.Prod.Name {
					t.Errorf("expected name %q got %q", expect.Prod.Name, new.Prod.Name)
				}
				if expect.Prod.Category != new.Prod.Category {
					t.Errorf("expected category %q got %q", expect.Prod.Category, new.Prod.Category)
				}
				if expect.Prod.Banned != new.Prod.Banned {
					t.Errorf("expected banned %t got %t", expect.Prod.Banned, new.Prod.Banned)
				}
				if expect.Prod.Preference != new.Prod.Preference {
					t.Errorf("expected preference %f got %f", expect.Prod.Preference, new.Prod.Preference)
				}
				if expect.Prod.SelectionScore != new.Prod.SelectionScore {
					t.Errorf("expected selection score %f got %f", expect.Prod.SelectionScore, new.Prod.SelectionScore)
				}

			}
		})
	}

}

func TestServiceProds_Add_Table(t *testing.T) {
	type testedProd struct {
		newProd      ProdsForGui
		expectedProd ProdsForGui
	}
	tests := []struct {
		name          string
		expectedError bool
		testProd      testedProd
	}{
		{
			name:          "add default",
			expectedError: false,
			testProd: testedProd{
				newProd: ProdsForGui{
					products.NewDefaultProduct("rice", products.Grain)},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         false,
						Preference:     products.Neutral,
						SelectionScore: 1.0,
					}},
			}},
		{
			name:          "add banned disliked",
			expectedError: false,
			testProd: testedProd{
				newProd: ProdsForGui{
					products.NewProduct(
						"rice", products.Grain, true, products.Disliked)},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "rice",
						Category:       products.Grain,
						Banned:         true,
						Preference:     products.Disliked,
						SelectionScore: 1.0,
					}},
			}},
		{
			name:          "add vegetable liked",
			expectedError: false,
			testProd: testedProd{
				newProd: ProdsForGui{
					products.NewProduct(
						"tomato", products.Vegetable, false, products.Liked)},
				expectedProd: ProdsForGui{
					products.Product{
						Name:           "tomato",
						Category:       products.Vegetable,
						Banned:         false,
						Preference:     products.Liked,
						SelectionScore: 1.0,
					}},
			}},
		{
			name:          "empty name",
			expectedError: true,
			testProd: testedProd{
				newProd: ProdsForGui{
					products.NewDefaultProduct("", products.Vegetable)},
			}},
		{
			name:          "empty category",
			expectedError: true,
			testProd: testedProd{
				newProd: ProdsForGui{
					products.NewDefaultProduct("rice", "")},
			}},
		{
			name:          "empty preference",
			expectedError: true,
			testProd: testedProd{
				newProd: ProdsForGui{
					products.NewProduct("rice", products.Grain, false, products.PreferenceStatus(0))},
			}},
		{
			name:          "empty all",
			expectedError: true,
			testProd: testedProd{
				newProd: ProdsForGui{
					products.NewProduct("", "", false, products.PreferenceStatus(0))},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var foundId products.ProductID = 1

			expect := tt.testProd.expectedProd
			db := setupDB(t)

			err := tt.testProd.newProd.Add(db)
			if (err != nil) != tt.expectedError {
				t.Fatalf("for error %v expected %t, got %t", err, tt.expectedError, (err != nil))
			}
			if !tt.expectedError {

				new := productFounder(t, db, foundId)

				if expect.Prod.Name != new.Prod.Name {
					t.Errorf("expected name %q got %q", expect.Prod.Name, new.Prod.Name)
				}
				if expect.Prod.Category != new.Prod.Category {
					t.Errorf("expected category %q got %q", expect.Prod.Category, new.Prod.Category)
				}
				if expect.Prod.Banned != new.Prod.Banned {
					t.Errorf("expected banned %t got %t", expect.Prod.Banned, new.Prod.Banned)
				}
				if expect.Prod.Preference != new.Prod.Preference {
					t.Errorf("expected preference %f got %f", expect.Prod.Preference, new.Prod.Preference)
				}
				if expect.Prod.SelectionScore != new.Prod.SelectionScore {
					t.Errorf("expected selection score %f got %f", expect.Prod.SelectionScore, new.Prod.SelectionScore)
				}

			}
		})
	}
}

func TestServiceProds_Delete(t *testing.T) {
	var foundID products.ProductID = 1
	db := setupDB(t)
	defer db.Close()

	p := ProdsForGui{products.NewDefaultProduct("rice", products.Grain)}
	insertDB(t, db, p.Prod)
	p = productFounder(t, db, foundID)
	err := p.Delete(db)
	if err != nil {
		t.Fatalf("expected %v got %v", nil, err)
	}
	_, err = storage.SelectProductByID(db, foundID)
	if err == nil {
		t.Fatalf("expected an error got %v", err)
	}
}

func TestGetList_GetList_Table(t *testing.T) {
	type setup struct {
		category   products.Category
		banned     bool
		preference products.PreferenceStatus
	}
	tests := []struct {
		name         string
		prod         []ProdsForGui
		expectedProd []ProdsForGui
		st           setup
		fn           func(*sql.DB, setup) ([]ProdsForGui, error)
	}{
		{
			name: "happy path get list",
			prod: []ProdsForGui{
				{products.NewDefaultProduct("rice", products.Grain)},
				{products.NewDefaultProduct("egg", products.Protein)},
				{products.NewDefaultProduct("tomato", products.Vegetable)}},
			expectedProd: []ProdsForGui{
				{products.Product{
					Name: "rice", Category: products.Grain,
					Banned: false, Preference: products.Neutral}},
				{products.Product{
					Name: "egg", Category: products.Protein,
					Banned: false, Preference: products.Neutral}},
				{products.Product{
					Name: "tomato", Category: products.Vegetable,
					Banned: false, Preference: products.Neutral}}},
			fn: func(db *sql.DB, s setup) ([]ProdsForGui, error) {
				return GetList(db)
			}},
		{
			name: "happy path get list filtered Grain",
			prod: []ProdsForGui{
				{products.NewProduct("rice", products.Grain, false, products.Neutral)},
				{products.NewProduct("egg", products.Protein, true, products.Liked)},
				{products.NewProduct("tomato", products.Vegetable, false, products.Liked)}},
			expectedProd: []ProdsForGui{
				{products.Product{
					Name: "rice", Category: products.Grain,
					Banned: false, Preference: products.Neutral}}},
			st: setup{category: products.Grain},
			fn: func(db *sql.DB, s setup) ([]ProdsForGui, error) {
				return GetListFiltered(db, &s.category, nil, nil)
			},
		},
		{

			name: "happy path get list filtered banned",
			prod: []ProdsForGui{
				{products.NewProduct("rice", products.Grain, false, products.Neutral)},
				{products.NewProduct("egg", products.Protein, true, products.Liked)},
				{products.NewProduct("tomato", products.Vegetable, false, products.Liked)},
			},
			expectedProd: []ProdsForGui{
				{products.Product{
					Name: "egg", Category: products.Protein,
					Banned: true, Preference: products.Liked}},
			},
			st: setup{banned: true},
			fn: func(db *sql.DB, s setup) ([]ProdsForGui, error) {
				return GetListFiltered(db, nil, &s.banned, nil)
			},
		},
		{
			name: "happy path get list filtered preference",
			prod: []ProdsForGui{
				{products.NewProduct("rice", products.Grain, false, products.Neutral)},
				{products.NewProduct("egg", products.Protein, true, products.Neutral)},
				{products.NewProduct("tomato", products.Vegetable, false, products.Liked)},
			},
			expectedProd: []ProdsForGui{
				{products.Product{
					Name: "tomato", Category: products.Vegetable,
					Banned: false, Preference: products.Liked}},
			},
			st: setup{preference: products.Liked},
			fn: func(db *sql.DB, s setup) ([]ProdsForGui, error) {
				return GetListFiltered(db, nil, nil, &s.preference)
			},
		},
		{
			name: "happy path get list filtered",
			prod: []ProdsForGui{
				{products.NewProduct("rice", products.Grain, false, products.Neutral)},
				{products.NewProduct("egg", products.Protein, true, products.Liked)},
				{products.NewProduct("tomato", products.Vegetable, false, products.Liked)},
			},
			expectedProd: []ProdsForGui{
				{products.Product{
					Name: "egg", Category: products.Protein,
					Banned: true, Preference: products.Liked}},
			},

			st: setup{category: products.Protein, banned: true, preference: products.Liked},
			fn: func(db *sql.DB, s setup) ([]ProdsForGui, error) {
				return GetListFiltered(db, &s.category, &s.banned, &s.preference)
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupDB(t)
			for _, elem := range tt.prod {
				insertDB(t, db, elem.Prod)
			}
			gotProd, err := tt.fn(db, tt.st)
			if err != nil {
				t.Fatalf("expected  error %v got %v", nil, err)
			}
			if len(gotProd) != len(tt.expectedProd) {
				t.Fatalf("unexpected count results data")
			}
			for i, elem := range gotProd {
				if elem.Prod.Name != tt.expectedProd[i].Prod.Name {
					t.Errorf("expected name %q got %q", tt.expectedProd[i].Prod.Name, elem.Prod.Name)
				}
				if elem.Prod.Category != tt.expectedProd[i].Prod.Category {
					t.Errorf("expected category %q got %q", tt.expectedProd[i].Prod.Category, elem.Prod.Category)
				}
				if elem.Prod.Banned != tt.expectedProd[i].Prod.Banned {
					t.Errorf("expected banned %t got %t", tt.expectedProd[i].Prod.Banned, elem.Prod.Banned)
				}
				if elem.Prod.Preference != tt.expectedProd[i].Prod.Preference {
					t.Errorf("expected %f got %f", tt.expectedProd[i].Prod.Preference, elem.Prod.Preference)
				}
			}
		})
	}
}

func TestDayControl_MinMaxToString_Table(t *testing.T) {
	tests := []struct {
		name     string
		min      int
		max      int
		expected []string
	}{{
		name:     "happy",
		min:      1,
		max:      6,
		expected: []string{"1", "2", "3", "4", "5", "6"},
	},
		{
			name:     "min more max",
			min:      6,
			max:      1,
			expected: []string{},
		},
		{
			name:     "min equal max",
			min:      6,
			max:      6,
			expected: []string{"6"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ConvertMinToMaxInString(tt.min, tt.max)
			if len(res) != len(tt.expected) {
				t.Fatalf("expected len %d got %d", len(tt.expected), len(res))
			}
			for i, elem := range tt.expected {
				if res[i] != elem {
					t.Errorf("expected %q got %q", elem, res[i])
				}
			}
		})
	}
}
