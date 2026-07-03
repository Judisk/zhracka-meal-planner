package gui

import (
	"database/sql"
	"fmt"
	"foods/internal/products"
	"foods/internal/service"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var sizes = []float32{50, 220, 100, 120, 80, 60, 60}

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

func productsTable(db *sql.DB, w fyne.Window, rightPanel *fyne.Container) (*widget.Table, error) {
	data, err := service.GetList(db)
	if err != nil {
		return nil, fmt.Errorf("create product table: %w", err)
	}
	var table *widget.Table

	table = widget.NewTable(
		func() (int, int) {
			return len(data) + 1, 7
		},
		func() fyne.CanvasObject {
			label := widget.NewLabel("")
			button := widget.NewButton("***", nil)

			return container.NewStack(label, button)
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			stack := cell.(*fyne.Container)

			label := stack.Objects[0].(*widget.Label)
			button := stack.Objects[1].(*widget.Button)

			label.Hide()
			button.Hide()

			if id.Row == 0 {
				label.Show()
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
				case 5:
					label.SetText("Editing")
				case 6:
					label.SetText("Deleting")
				}

				return
			}
			prod := data[id.Row-1]
			switch id.Col {
			case 0:
				label.Show()
				label.SetText(strconv.Itoa(int(prod.Prod.ID)))
			case 1:
				label.Show()
				label.SetText(prod.Prod.Name)
			case 2:
				label.Show()
				label.SetText(string(prod.Prod.Category))
			case 3:
				label.Show()
				if prod.Prod.Banned {
					label.SetText("Yes")
				} else {
					label.SetText("No")
				}

			case 4:
				label.Show()
				switch prod.Prod.Preference {
				case products.Neutral:
					label.SetText("Neutral")
				case products.Liked:
					label.SetText("Liked")
				case products.Disliked:
					label.SetText("Disliked")
				}
			case 5:
				button.Show()
				button.SetText("Edit")
				button.Importance = widget.WarningImportance
				button.OnTapped = func() { EditingButton(db, w, rightPanel, prod) }
			case 6:
				button.Show()
				button.SetText("Delete")
				button.Importance = widget.DangerImportance
				button.OnTapped = func() { DeleteConfirmButton(db, w, rightPanel, prod) }
			}
		})

	return defaultSizeTable(table), nil
}

func DeleteConfirmButton(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, p service.ProdsForGui) {

	dialog.ShowConfirm("Confirm", "Confirm Deleting?", func(b bool) {
		if !b {
			return
		}
		if err := p.Delete(db); err != nil {
			dialog.ShowError(err, w)
		}
		newTable, err := productsTable(db, w, rightPanel)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rightPanel.Objects[0] = newTable
		rightPanel.Refresh()
	}, w)

}

func EditingButton(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, prod service.ProdsForGui) {
	p := prod.Prod

	nameEntry := widget.NewEntry()
	nameEntry.SetText(p.Name)

	categorySelect := widget.NewSelect([]string{"grain", "protein", "vegetable"}, nil)
	categorySelect.SetSelected(string(p.Category))

	textPref := ""
	switch p.Preference {
	case products.Liked:
		textPref = "Liked"
	case products.Neutral:
		textPref = "Neutral"
	case products.Disliked:
		textPref = "Disliked"
	}
	preferenceSelect := widget.NewSelect([]string{"Liked", "Neutral", "Disliked"}, nil)
	preferenceSelect.SetSelected(textPref)

	bannedCheck := widget.NewCheck("Banned", nil)
	bannedCheck.SetChecked(p.Banned)

	saveButton := widget.NewButton("Save", func() {
		var pref products.PreferenceStatus
		switch preferenceSelect.Selected {
		case "Liked":
			pref = products.Liked
		case "Neutral":
			pref = products.Neutral
		case "Disliked":
			pref = products.Disliked
		}
		newProd := service.ProdsForGui{
			Prod: products.Product{
				ID:         p.ID,
				Name:       nameEntry.Text,
				Category:   products.Category(categorySelect.Selected),
				Banned:     bannedCheck.Checked,
				Preference: products.PreferenceStatus(pref),
			},
		}
		if err := newProd.Edit(db); err != nil {
			dialog.ShowError(err, w)
			return
		}
		newTable, err := productsTable(db, w, rightPanel)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rightPanel.Objects[0] = newTable
		rightPanel.Refresh()
	})

	fields := container.NewHBox(categorySelect, preferenceSelect, bannedCheck)
	content := container.NewVBox(nameEntry, fields, saveButton)

	dialog.ShowCustom("Editing", "Close", content, w)

}
