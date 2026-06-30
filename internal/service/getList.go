package service

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"foods/internal/storage"
)

func GetList(db *sql.DB) ([]products.Product, error) {
	p, err := storage.SelectAll(db)
	if err != nil {
		return nil, fmt.Errorf("get product list: %w", err)
	}
	return p, nil
}
