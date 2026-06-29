package gui

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"foods/internal/service"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var sizes = []float32{50, 250, 100, 120, 80}

func sumOfSize() float32 {
	var result float32
	for _, elem := range sizes {
		result += elem
	}
	return result
}

func defaultSizeTable(table *widget.Table) *widget.Table {
	for i, elem := range sizes {
		table.SetColumnWidth(i, float32(elem))
	}
	return table
}

func productsTable(db *sql.DB) (*widget.Table, error) {
	data, err := service.GetList(db)
	if err != nil {
		return nil, fmt.Errorf("create list: %w", err)
	}
	table := widget.NewTable(
		func() (int, int) {
			return len(data) + 1, 5
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {

			label := cell.(*widget.Label)

			if id.Row == 0 {
				switch id.Col {
				case 0:
					label.SetText("ID")
				case 1:
					label.SetText("Name")
				case 2:
					label.SetText("Category")
				case 3:
					label.SetText("Banned")

				case 4:
					label.SetText("Preference")
				}
				return
			}
			prod := data[id.Row-1]
			switch id.Col {
			case 0:
				label.SetText(strconv.Itoa(int(prod.ID)))
			case 1:
				label.SetText(prod.Name)
			case 2:
				label.SetText(string(prod.Category))
			case 3:
				if prod.Banned {
					label.SetText("Да")
				} else {
					label.SetText("Нет")
				}

			case 4:
				switch prod.Preference {
				case products.Neutral:
					label.SetText("Нейтральный")
				case products.Liked:
					label.SetText("Любимый")
				case products.Disliked:
					label.SetText("Нелюбимые")
				}
			}
		})

	return defaultSizeTable(table), nil
}
