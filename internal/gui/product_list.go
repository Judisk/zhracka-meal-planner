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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FiltredState struct {
	CategoryState    *products.Category
	BannedState      *bool
	PreferencesState *products.PreferenceStatus
}

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

func tableContainer(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, state FiltredState) (*fyne.Container, error) {
	table, err := productsTable(db, w, rightPanel, state)
	if err != nil {
		return nil, fmt.Errorf("create product table: %w", err)
	}
	addButton := widget.NewButton(
		"ADD", func() {})
	addButton.Importance = widget.HighImportance
	addButton.OnTapped = func() {
		AddButton(db, w, rightPanel, state)

	}
	header := headerFn()

	return container.NewBorder(container.NewVBox(addButton, header), nil, nil, nil, table), nil
}

func productsTable(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, state FiltredState) (*widget.Table, error) {
	data, err := service.GetListFiltered(db, state.CategoryState, state.BannedState, state.PreferencesState)
	if err != nil {
		return nil, fmt.Errorf("create product table: %w", err)
	}
	var table *widget.Table

	table = widget.NewTable(
		func() (int, int) {
			return len(data), 7
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

			prod := data[id.Row]
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
				button.Refresh()
				button.OnTapped = func() { EditingButton(db, w, rightPanel, prod, state) }
			case 6:
				button.Show()
				button.SetText("Delete")
				button.Importance = widget.DangerImportance
				button.Refresh()
				button.OnTapped = func() { DeleteConfirmButton(db, w, rightPanel, prod, state) }
			}
		})

	return defaultSizeTable(table), nil
}

func DeleteConfirmButton(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, p service.ProdsForGui, state FiltredState) {

	dialog.ShowConfirm("Confirm", "Confirm Deleting?", func(b bool) {
		if !b {
			return
		}
		if err := p.Delete(db); err != nil {
			dialog.ShowError(err, w)
		}
		newTable, err := tableContainer(db, w, rightPanel, state)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rightPanel.Objects[0] = newTable
		rightPanel.Refresh()
	}, w)

}

func EditingButton(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, prod service.ProdsForGui, state FiltredState) {
	p := prod.Prod

	nameEntry := widget.NewEntry()
	nameEntry.SetText(p.Name)

	categorySelect := widget.NewSelect([]string{"Grain", "Protein", "Vegetable"}, nil)
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
		newTable, err := tableContainer(db, w, rightPanel, state)
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

func AddButton(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, state FiltredState) {
	nameEntry := widget.NewEntry()
	nameEntry.PlaceHolder = "Enter product name"

	categorySelect := widget.NewSelect([]string{"Grain", "Protein", "Vegetable"}, nil)

	preferenceSelect := widget.NewSelect([]string{"Liked", "Neutral", "Disliked"}, nil)

	bannedDefault := false
	bannedCheck := widget.NewCheck("Banned", nil)
	bannedCheck.SetChecked(bannedDefault)

	var d dialog.Dialog
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
				Name:       nameEntry.Text,
				Category:   products.Category(categorySelect.Selected),
				Banned:     bannedCheck.Checked,
				Preference: products.PreferenceStatus(pref),
			},
		}
		if err := newProd.Add(db); err != nil {
			dialog.ShowError(err, w)
			return
		}
		newTable, err := tableContainer(db, w, rightPanel, state)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rightPanel.Objects[0] = newTable
		rightPanel.Refresh()
		d.Hide()
	})
	fields := container.NewHBox(categorySelect, preferenceSelect, bannedCheck)
	content := container.NewVBox(nameEntry, fields, saveButton)
	d = dialog.NewCustom("Add", "Close", content, w)
	d.Show()

}

func headerFn() *fyne.Container {
	headers := []string{"ID", "Name", "Category", "Banned", "Preference"}
	widths := sizes[:5]
	items := []fyne.CanvasObject{}
	for i, text := range headers {
		w := widths[i]
		btn := widget.NewButton(text, func() {})
		items = append(items, container.New(layout.NewGridWrapLayout(fyne.NewSize(w, 35)), btn))
	}
	return container.NewHBox(items...)
}

func CategoryFilter(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, state FiltredState) *widget.Select {
	s := widget.NewSelect([]string{"Category", "Grain", "Protein", "Vegetable"}, func(selected string) {
		newState := state
		if selected == "All" {
			newState.CategoryState = nil
		} else {
			cat := products.Category(selected)
			newState.CategoryState = &cat
		}
		newContainer, err := tableContainer(db, w, rightPanel, state)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rightPanel.Objects[0] = newContainer
		rightPanel.Refresh()
	})
	s.SetSelected("All")
	return s

}

func preferenceFilter(db *sql.DB, w fyne.Window, rightPanel *fyne.Container, state FiltredState) *widget.Select {
	s := widget.NewSelect([]string{"Preference", "Liked", "Neutral", "Disliked"}, func(selected string) {
		newState := state
		if selected == "Preference" {
			newState.PreferencesState = nil
		} else {
			var pref products.PreferenceStatus
			switch selected {
			case "Liked":
				pref = products.Liked
			case "Neutral":
				pref = products.Neutral
			case "Disliked":
				pref = products.Disliked
			}

			newState.PreferencesState = &pref
		}
		newContainer, err := tableContainer(db, w, rightPanel, state)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rightPanel.Objects[0] = newContainer
		rightPanel.Refresh()
	})
	s.SetSelected("Preference")
	return s

}
