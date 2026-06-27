package storage

import (
	"database/sql"
	"foods/internal/products"
	"strings"
	"testing"
)

// Helpers
func setupDB(t *testing.T) *sql.DB {
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return db
}

func insertDB(t *testing.T, db *sql.DB, p products.Product) {
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

func TestDB_SelectSuccess(t *testing.T) {
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

func TestDB_SelectAllowedProductsSuccess(t *testing.T) {

	productsNames := []string{"T1", "T2", "T3", "B1"}
	expectedIDs := []products.ProductID{1, 2, 3}
	category := products.Grain

	expectedLen := len(productsNames) - 1
	expectedNotBanned := false
	expectedPreference := products.Neutral
	expectedSelectionScore := 1.0

	allowedStatus := false
	bannedStatus := true

	db := setupDB(t)
	defer db.Close()

	for i := range expectedLen {
		insertDB(t, db, products.NewProduct(productsNames[i], category, allowedStatus, products.Neutral))
	}

	insertDB(t, db, products.NewProduct(productsNames[3], category, bannedStatus, products.Neutral))

	gotAllowedProducts, err := SelectAllowedProductsByCategory(db, category)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(gotAllowedProducts) != expectedLen {
		t.Fatalf("expected %d products, got %d", expectedLen, len(gotAllowedProducts))
	}

	for i, elem := range gotAllowedProducts {
		if elem.ID != expectedIDs[i] {
			t.Errorf("expected ID %d, got %d", expectedIDs[i], elem.ID)
		}

		if elem.Name != productsNames[i] {
			t.Errorf("expected name %q, got %q", productsNames[i], elem.Name)
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

}

func TestDB_SelectAllProductsSuccess(t *testing.T) {

	productsNames := []string{"T1", "T2", "T3"}
	expectedIDs := []products.ProductID{1, 2, 3}
	category := products.Grain
	expectedBanned := false
	expectedPreference := products.Neutral
	expectedSelectionScore := 1.0

	db := setupDB(t)
	defer db.Close()

	for _, name := range productsNames {
		insertDB(t, db, products.NewDefaultProduct(name, category))
	}

	gotProducts, err := SelectAllProductsByCategory(db, category)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(gotProducts) != len(productsNames) {
		t.Fatalf("expected %d products, got %d", len(productsNames), len(gotProducts))
	}

	for i, elem := range gotProducts {
		if elem.ID != expectedIDs[i] {
			t.Errorf("expected ID %d, got %d", expectedIDs[i], elem.ID)
		}

		if elem.Name != productsNames[i] {
			t.Errorf("expected name %q, got %q", productsNames[i], elem.Name)
		}

		if elem.Category != category {
			t.Errorf("expected category %q, got %q", category, elem.Category)
		}

		if elem.Banned != expectedBanned {
			t.Errorf("expected banned %t, got %t", expectedBanned, elem.Banned)
		}

		if elem.Preference != expectedPreference {
			t.Errorf("expected preference %f, got %f", expectedPreference, elem.Preference)
		}
		if elem.SelectionScore != expectedSelectionScore {
			t.Errorf("expected selection_score %f, got %f", expectedSelectionScore, elem.SelectionScore)
		}
	}
}

func TestDB_SelectBannedProductsCheckIDandNameSuccess(t *testing.T) {
	allowedProduct := products.NewDefaultProduct("A1", products.Grain)
	bannedProduct := products.NewProduct("B1", products.Grain, true, products.Neutral)

	var expectedBannedID products.ProductID = 2
	var expectedBannedName string = "b1"

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

	if !gotBannedProduct.ByID[expectedBannedID] {
		t.Errorf("expected banned product id %d to be present in ByID", expectedBannedID)
	}

	if !gotBannedProduct.ByName[expectedBannedName] {
		t.Errorf("expected banned product name %q to be present in ByName", expectedBannedName)
	}

	if gotBannedProduct.ByID[unexpectedID] {
		t.Errorf("unexpected allowed product id %d to be present in ByID", unexpectedID)
	}

	if gotBannedProduct.ByName[unexpectedName] {
		t.Errorf("unexpected allowed product name %q to be present in ByName", unexpectedName)
	}
}

func TestDB_SelectBannedProductsCheckFalseIDandTrueNameSuccess(t *testing.T) {
	allowedProduct := products.NewDefaultProduct("A1", products.Grain)
	bannedProduct := products.NewProduct("B1", products.Grain, true, products.Neutral)

	var unexpectedBannedID products.ProductID = 22
	var expectedBannedName string = "b1"

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

	if gotBannedProduct.ByID[unexpectedBannedID] {
		t.Errorf("expected unknown product id %d to be absent from ByID", unexpectedBannedID)
	}

	if !gotBannedProduct.ByName[expectedBannedName] {
		t.Errorf("expected banned product name %q to be present in ByName", expectedBannedName)
	}

	if gotBannedProduct.ByID[unexpectedID] {
		t.Errorf("unexpected allowed product id %d to be present in ByID", unexpectedID)
	}

	if gotBannedProduct.ByName[unexpectedName] {
		t.Errorf("unexpected allowed product name %q to be present in ByName", unexpectedName)
	}
}

func TestDB_SelectBannedProductsCheckTrueIDandFalseNameSuccess(t *testing.T) {
	allowedProduct := products.NewDefaultProduct("A1", products.Grain)
	bannedProduct := products.NewProduct("B1", products.Grain, true, products.Neutral)

	var expectedBannedID products.ProductID = 2
	var unexpectedBannedName string = "b11"

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

	if !gotBannedProduct.ByID[expectedBannedID] {
		t.Errorf("expected banned product id %d to be present in ByID", expectedBannedID)
	}

	if gotBannedProduct.ByName[unexpectedBannedName] {
		t.Errorf("expected unknown product name %q to be absent from ByName", unexpectedBannedName)
	}

	if gotBannedProduct.ByID[unexpectedID] {
		t.Errorf("unexpected allowed product id %d to be present in ByID", unexpectedID)
	}

	if gotBannedProduct.ByName[unexpectedName] {
		t.Errorf("unexpected allowed product name %q to be present in ByName", unexpectedName)
	}
}

// Errors path ///////////////////////////////////////////////
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
