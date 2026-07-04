package service

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"foods/internal/storage"
)

func GetList(db *sql.DB) ([]ProdsForGui, error) {
	prods := []ProdsForGui{}
	p, err := storage.SelectAll(db)
	if err != nil {
		return nil, fmt.Errorf("get product list: %w", err)
	}
	for _, elem := range p {
		prods = append(prods, ProdsForGui{Prod: elem})
	}
	return prods, nil
}

func GetListFiltered(db *sql.DB, category *products.Category, banned *bool, preference *products.PreferenceStatus) ([]ProdsForGui, error) {
	prods := []ProdsForGui{}
	p, err := storage.SelectAllFiltered(db, category, banned, preference)
	if err != nil {
		return nil, fmt.Errorf("get product list: %w", err)
	}
	for _, elem := range p {
		prods = append(prods, ProdsForGui{Prod: elem})
	}
	return prods, nil
}
