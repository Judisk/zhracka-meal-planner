package gui

import (
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func optionsList(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, state FiltredState) *fyne.Container {
	return container.NewVBox(
		allListButton(rightPanel, db, w, state),
		oneDayButton(rightPanel, db),
		oneDishButton(rightPanel, db),
	)
}

func allListButton(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, state FiltredState) *widget.Button {
	return widget.NewButton("All Products", func() {
		productTableContainer, err := tableContainer(db, w, rightPanel, state)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		rightPanel.Objects[0] = productTableContainer
		rightPanel.Refresh()
	})
}
func oneDayButton(rightPanel *fyne.Container, db *sql.DB) *widget.Button {
	return widget.NewButton("One Day", func() {
		//rightPanel.Objects[0] = dayView
		rightPanel.Refresh()
	})
}

func oneDishButton(rightPanel *fyne.Container, db *sql.DB) *widget.Button {
	return widget.NewButton("One Dish", func() {
		rightPanel.Refresh()
	})
}
