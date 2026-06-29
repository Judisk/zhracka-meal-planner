package storage

import (
	"database/sql"
	"foods/internal/products"
	"strings"
	"testing"
)

// Helpers
func setupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return db
}

func insertDB(t *testing.T, db *sql.DB, p products.Product) {
	t.Helper()
	err := InsertProduct(db, p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// HAPPY PATH /////////////////////////////////////////////////
func TestDB_CreatingSuccess(t *testing.T) {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()
}

func TestDB_InsertSuccess(t *testing.T) {
	name := "Test"
	category := products.Grain

	db := setupDB(t)
	defer db.Close()

	if err := InsertProduct(db, products.NewDefaultProduct(name, category)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDB_DeleteSuccess(t *testing.T) {
	name := "Test"
	category := products.Grain
	var testingID products.ProductID = 1

	db := setupDB(t)
	defer db.Close()

	insertDB(t, db, products.NewDefaultProduct(name, category))
	if err := DeleteProductsByName(db, name); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := SelectProductByID(db, testingID)
	if err == nil {
		t.Fatal("expected error (product not found), but got nil")
	}

}

func TestDB_SelectIDsSuccess(t *testing.T) {
	name1 := "Test1"
	name2 := "Test2"
	expectedID := []products.ProductID{1, 2}
	var testingID products.ProductID = 1
	category := products.Grain
	expectedLen := 2

	expectedBanned := false
	expectedPreference := products.Neutral

	db := setupDB(t)
	defer db.Close()

	insertDB(t, db, products.NewDefaultProduct(name1, category))

	insertDB(t, db, products.NewDefaultProduct(name2, category))

	arr, err := SelectProductIDsByCategory(db, category)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(arr) != expectedLen {
		t.Fatalf("expected ID slice length %d, got %d", expectedLen, len(arr))
	}

	for i := range arr {
		if arr[i] != expectedID[i] {
			t.Errorf("expected ID %d, got %d", expectedID[i], arr[i])
		}
	}

	p, err := SelectProductByID(db, testingID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.ID != testingID {
		t.Errorf("expected ID %d, got %d", testingID, p.ID)
	}

	if p.Name != name1 {
		t.Errorf("expected name %q, got %q", name1, p.Name)
	}

	if p.Category != category {
		t.Errorf("expected category %q, got %q", category, p.Category)
	}

	if p.Banned != expectedBanned {
		t.Errorf("expected banned %t, got %t", expectedBanned, p.Banned)
	}

	if p.Preference != expectedPreference {
		t.Errorf("expected preference %f, got %f", expectedPreference, p.Preference)
	}
}

func TestDB_SelectsSuccess(t *testing.T) {

	tests := []struct {
		name        string
		testNames   []string
		queryFn     func(*sql.DB, products.Category) ([]products.Product, error)
		expectedLen int
		needAdd     bool
	}{{
		name:        "Allowed Products",
		testNames:   []string{"T1", "T2", "T3", "B1"},
		queryFn:     SelectUnbannedProductsByCategory,
		expectedLen: 3,
		needAdd:     true,
	}, {
		name:        "All products",
		testNames:   []string{"T1", "T2", "T3"},
		queryFn:     SelectAllProductsByCategory,
		expectedLen: 3,
		needAdd:     false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedIDs := []products.ProductID{1, 2, 3}
			category := products.Grain
			expectedNotBanned := false
			expectedPreference := products.Neutral
			expectedSelectionScore := 1.0

			db := setupDB(t)
			defer db.Close()

			for i := range tt.expectedLen {
				insertDB(t, db, products.NewDefaultProduct(tt.testNames[i], category))
			}
			if tt.needAdd {
				insertDB(t, db, products.NewProduct(tt.testNames[tt.expectedLen], category, true, products.Neutral))
			}
			gotData, err := tt.queryFn(db, category)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(gotData) != tt.expectedLen {
				t.Fatalf("expected %d products, got %d", tt.expectedLen, len(gotData))
			}

			for i, elem := range gotData {
				if elem.ID != expectedIDs[i] {
					t.Errorf("expected ID %d, got %d", expectedIDs[i], elem.ID)
				}

				if elem.Name != tt.testNames[i] {
					t.Errorf("expected name %q, got %q", tt.testNames[i], elem.Name)
				}

				if elem.Category != category {
					t.Errorf("expected category %q, got %q", category, elem.Category)
				}

				if elem.Banned != expectedNotBanned {
					t.Errorf("expected banned %t, got %t", expectedNotBanned, elem.Banned)
				}

				if elem.Preference != expectedPreference {
					t.Errorf("expected preference %f, got %f", expectedPreference, elem.Preference)
				}
				if elem.SelectionScore != expectedSelectionScore {
					t.Errorf("expected selection_score %f, got %f", expectedSelectionScore, elem.SelectionScore)
				}
			}
		})
	}
}

func TestDB_Table_SelectBannedProducts(t *testing.T) {

	tests := []struct {
		name string

		expectedBannedID   products.ProductID
		expectedBannedName string

		idShouldExist   bool
		nameShouldExist bool
	}{
		{name: "SelectBannedProductsCheckIDandNameSuccess",

			expectedBannedID:   2,
			expectedBannedName: "b1",

			idShouldExist:   true,
			nameShouldExist: true,
		}, {
			name: "SelectBannedProductsCheckFalseIDandTrueNameSuccess",

			expectedBannedID:   22,
			expectedBannedName: "b1",

			idShouldExist:   false,
			nameShouldExist: true,
		}, {
			name: "SelectBannedProductsCheckTrueIDandFalseNameSuccess",

			expectedBannedID:   2,
			expectedBannedName: "b11",

			idShouldExist:   true,
			nameShouldExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			allowedProduct := products.NewDefaultProduct("A1", products.Grain)
			bannedProduct := products.NewProduct("B1", products.Grain, true, products.Neutral)

			var unexpectedID products.ProductID = 1
			var unexpectedName string = "a1"

			db := setupDB(t)
			defer db.Close()

			insertDB(t, db, allowedProduct)
			insertDB(t, db, bannedProduct)

			gotBannedProduct, err := SelectBannedProducts(db)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotBannedProduct.ByID[tt.expectedBannedID] != tt.idShouldExist {
				t.Errorf("expected banned product id %d to be present in ByID", tt.expectedBannedID)
			}

			if gotBannedProduct.ByName[tt.expectedBannedName] != tt.nameShouldExist {
				t.Errorf("expected banned product name %q to be present in ByName", tt.expectedBannedName)
			}

			if gotBannedProduct.ByID[unexpectedID] {
				t.Errorf("unexpected allowed product id %d to be present in ByID", unexpectedID)
			}

			if gotBannedProduct.ByName[unexpectedName] {
				t.Errorf("unexpected allowed product name %q to be present in ByName", unexpectedName)
			}

		})
	}
}

func TestDB_UpdateSelectionScoreSuccess(t *testing.T) {
	Prods := []products.Product{
		products.NewDefaultProduct("A1", products.Grain),
		products.NewProduct("A2", products.Grain, false, products.Liked),
		products.NewDefaultProduct("R1", products.Grain),
	}

	expectedScoreProds := []float64{
		2.0,
		2.5,
		0.0}

	var Prod3SQLID products.ProductID = 3

	db := setupDB(t)
	defer db.Close()
	for _, elem := range Prods {
		insertDB(t, db, elem)
	}

	prod3, err := SelectProductByID(db, Prod3SQLID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ManyResets(db, prod3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := UpdateSelectionScore(db); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prods, err := SelectAllProductsByCategory(db, products.Grain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, elem := range prods {
		if elem.SelectionScore != expectedScoreProds[i] {
			t.Errorf("expected score %f got %f", expectedScoreProds[i], elem.SelectionScore)
		}
	}
}

func TestDB_ResetScoreSuccess(t *testing.T) {
	Prod1 := products.NewDefaultProduct("A1", products.Grain)
	expectedScoreProd := -1.0

	db := setupDB(t)
	defer db.Close()

	insertDB(t, db, Prod1)
	prod1, err := SelectAllProductsByCategory(db, products.Grain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := ManyResets(db, prod1...); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	prods, err := SelectAllProductsByCategory(db, products.Grain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if prods[0].SelectionScore != expectedScoreProd {
		t.Errorf("expected score %f got %f", expectedScoreProd, prods[0].SelectionScore)
	}
}

// Errors path ///////////////////////////////////////////////

func TestDb_UpdateSelectionScoreClosedDb(t *testing.T) {
	Prod1 := products.NewDefaultProduct("A1", products.Grain)
	db := setupDB(t)
	insertDB(t, db, Prod1)

	db.Close()

	if err := UpdateSelectionScore(db); err == nil {
		t.Fatalf("expected an error: %v", err)
	}

}

func TestDb_ResetScoreClosedDb(t *testing.T) {
	Prod1 := products.NewDefaultProduct("A1", products.Grain)
	db := setupDB(t)
	insertDB(t, db, Prod1)

	prod1, err := SelectAllProductsByCategory(db, products.Grain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	db.Close()

	if err := ManyResets(db, prod1...); err == nil {
		t.Fatalf("expected an error: %v", err)
	}
}

func TestDB_ResetScoreEmptyProds(t *testing.T) {
	db := setupDB(t)
	defer db.Close()

	prod1, err := SelectAllProductsByCategory(db, products.Grain)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := ManyResets(db, prod1...); err == nil {
		t.Fatalf("expected an error: %v", err)
	}
}

func TestDB_InsertEmptyName(t *testing.T) {
	name := ""
	category := products.Grain

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	err = InsertProduct(db, products.NewDefaultProduct(name, category))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDB_SelectBannedProductsReturnsEmptyWhenNoBannedProducts(t *testing.T) {

	allowedProduct := products.NewDefaultProduct("A1", products.Grain)

	db := setupDB(t)
	defer db.Close()

	insertDB(t, db, allowedProduct)

	gotBannedProduct, err := SelectBannedProducts(db)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotBannedProduct.ByName) != 0 {
		t.Errorf("expected empty ByName map, got %d items", len(gotBannedProduct.ByName))
	}
	if len(gotBannedProduct.ByID) != 0 {
		t.Errorf("expected empty ByID map, got %d items", len(gotBannedProduct.ByID))
	}

}

func TestDB_InsertClosedDB(t *testing.T) {
	name := "Test Name"
	category := products.Grain
	expectedError := "database is closed"

	db := setupDB(t)
	db.Close()

	err := InsertProduct(db, products.NewDefaultProduct(name, category))
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error to contain %q, got %v", expectedError, err)
	}
}

func TestDB_SelectUnknownCategoryReturnsEmptyIDs(t *testing.T) {
	unknownCategory := products.Category("testCat")
	name := "Test Name"
	category := products.Grain

	db := setupDB(t)
	defer db.Close()

	insertDB(t, db, products.NewDefaultProduct(name, category))

	arr, err := SelectProductIDsByCategory(db, unknownCategory)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(arr) != 0 {
		t.Errorf("expected empty ID slice, got %v", arr)
	}
}

func TestDB_SelectIDNotExist(t *testing.T) {
	var unknownID products.ProductID = 15

	db := setupDB(t)
	defer db.Close()

	_, err := SelectProductByID(db, unknownID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDB_DeleteNonExistentProduct(t *testing.T) {
	name := "Test"
	category := products.Grain
	removedName := "ThisProductDoesNotExist"
	var testingID products.ProductID = 1

	db := setupDB(t)
	defer db.Close()

	insertDB(t, db, products.NewDefaultProduct(name, category))

	if err := DeleteProductsByName(db, removedName); err != nil {
		t.Fatalf("expected NO error when deleting non-existent product, got: %v", err)
	}
	p, err := SelectProductByID(db, testingID)
	if err != nil {
		t.Fatalf("product 'Test' should still exist, got error: %v", err)
	}
	if p.Name != name {
		t.Errorf("product 'Test' was unexpectedly deleted or modified, got name: '%s'", p.Name)
	}
}

func TestDB_DeleteProductsByNameClosedDB(t *testing.T) {
	category := products.Grain
	name := "Test"

	db := setupDB(t)

	insertDB(t, db, products.NewDefaultProduct(name, category))

	db.Close()

	if err := DeleteProductsByName(db, name); err == nil {
		t.Fatal("expected error, got nil")
	}
}
