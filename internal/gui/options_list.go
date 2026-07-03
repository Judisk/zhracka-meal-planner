package gui

import (
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func optionsList(rightPanel *fyne.Container, db *sql.DB, w fyne.Window) *fyne.Container {
	return container.NewVBox(
		allListButton(rightPanel, db, w),
		oneDayButton(rightPanel, db),
		oneDishButton(rightPanel, db),
	)
}

func allListButton(rightPanel *fyne.Container, db *sql.DB, w fyne.Window) *widget.Button {
	return widget.NewButton("All Products", func() {
		productTableContainer, err := productsTable(db, w, rightPanel)
		if err != nil {
			// TODO: handle the error in the GUI
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
