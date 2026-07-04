package service

import (
	"database/sql"
	"errors"
	"fmt"
	"foods/internal/products"
	"foods/internal/storage"
)

type ProdsForGui struct {
	Prod products.Product
}

var ErrEmptyArgsForEditing = errors.New("zero arguments for editing")

func (p ProdsForGui) Edit(db *sql.DB) error {

	if err := storage.UpdateProductInfo(db, p.Prod); err != nil {
		return fmt.Errorf("edit product: %w", err)
	}
	return nil
}

func (p ProdsForGui) Delete(db *sql.DB) error {
	if err := storage.DeleteProductsByID(db, p.Prod.ID); err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	return nil
}

func (p ProdsForGui) Add(db *sql.DB) error {
	return storage.InsertProduct(db, p.Prod)
}
