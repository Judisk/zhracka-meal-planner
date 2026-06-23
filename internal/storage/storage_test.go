package storage

import (
	"foods/internal/products"
	"strings"
	"testing"
)

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

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	err = InsertProduct(db, name, category, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDB_DeleteSuccess(t *testing.T) {
	name := "Test"
	category := products.Grain
	var testingID products.ProductID = 1

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	if err = InsertProduct(db, name, category, false, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err = DeleteProductsByName(db, name); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = SelectProductByID(db, testingID)
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
	expectedFavorite := false

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	err = InsertProduct(db, name1, category, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = InsertProduct(db, name2, category, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

	if p.Favorite != expectedFavorite {
		t.Errorf("expected favorite %t, got %t", expectedFavorite, p.Favorite)
	}
}

func TestDB_SelectProductsSuccess(t *testing.T) {

	productsNames := []string{"T1", "T2", "T3"}
	expectedIDs := []products.ProductID{1, 2, 3}
	category := products.Grain
	expectedBanned := false
	expectedFavorite := false

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	for _, elem := range productsNames {
		if err = InsertProduct(db, elem, category, false, false); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	gotProducts, err := SelectProductsByCategory(db, category)
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

		if elem.Favorite != expectedFavorite {
			t.Errorf("expected favorite %t, got %t", expectedFavorite, elem.Favorite)
		}
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

	err = InsertProduct(db, name, category, false, false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDB_InsertClosedDB(t *testing.T) {
	name := "Test Name"
	category := products.Grain
	expectedError := "database is closed"

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	db.Close()

	err = InsertProduct(db, name, category, false, false)
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

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	err = InsertProduct(db, name, category, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

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

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	_, err = SelectProductByID(db, unknownID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDB_DeleteNonExistentProduct(t *testing.T) {
	name := "Test"
	category := products.Grain
	removedName := "ThisProductDoesNotExist"
	var testingID products.ProductID = 1

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()

	if err = InsertProduct(db, name, category, false, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err = DeleteProductsByName(db, removedName); err != nil {
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

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err = InsertProduct(db, name, category, false, false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	db.Close()

	if err = DeleteProductsByName(db, name); err == nil {
		t.Fatal("expected error, got nil")
	}
}
