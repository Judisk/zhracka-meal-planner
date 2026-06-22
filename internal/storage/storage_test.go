package storage

import (
	"foods/internal/products"
	"testing"
)

func TestDB_CreatingSuccess(t *testing.T) {
	_, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

}

func TestDB_InsertSuccess(t *testing.T) {
	name := "Test"
	category := products.Grain

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

	defer db.Close()

	err = InsertProduct(db, name, category, false, false)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

}

func TestDB_SelectSuccess(t *testing.T) {
	name1 := "Test1"
	name2 := "Test2"
	expectedId := []products.ProductID{1, 2}
	var testingId products.ProductID = 1
	category := products.Grain
	expectedLen := 2

	expectedBanned := false
	expectedFavorite := false

	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

	defer db.Close()

	err = InsertProduct(db, name1, category, false, false)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

	err = InsertProduct(db, name2, category, false, false)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}

	arr, err := SelectProductIDsByCategory(db, category)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}
	if len(arr) != expectedLen {
		t.Fatalf("ошибка длины массива: %v", err)
	}

	for i := range arr {
		if arr[i] != expectedId[i] {
			t.Errorf("ожидали '%d' получили %d", expectedId[i], arr[i])
		}
	}
	p, err := SelectProductByID(db, testingId)
	if err != nil {
		t.Fatalf("unexpected error : %v", err)
	}
	if p.ID != products.ProductID(testingId) {
		t.Errorf("ожидали '%d' получили %d", testingId, p.ID)
	}
	if p.Name != name1 {
		t.Errorf("ожидали '%s' получили %s", name1, p.Name)
	}
	if p.Category != category {
		t.Errorf("ожидали '%v' получили %v", category, p.Category)
	}
	if p.Banned != expectedBanned {
		t.Errorf("ожидали '%t' получили %t", expectedBanned, p.Banned)
	}
	if p.Favorite != expectedFavorite {
		t.Errorf("ожидали '%t' получили %t", expectedFavorite, p.Favorite)
	}

}
