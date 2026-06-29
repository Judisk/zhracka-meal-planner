package gui

import (
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func optionsList(rightPanel *fyne.Container, db *sql.DB) *fyne.Container {
	return container.NewVBox(
		allListButton(rightPanel, db),
		oneDayButton(rightPanel, db),
		oneDishButton(rightPanel, db),
	)
}

func allListButton(rightPanel *fyne.Container, db *sql.DB) *widget.Button {
	return widget.NewButton("AllList", func() {
		productTableContainer, err := productsTable(db)
		if err != nil {
			//обработать ошибку гуишным способом
		}
		rightPanel.Objects[0] = productTableContainer
		rightPanel.Refresh()
	})
}
func oneDayButton(rightPanel *fyne.Container, db *sql.DB) *widget.Button {
	return widget.NewButton("1 day", func() {
		//rightPanel.Objects[0] = dayView
		rightPanel.Refresh()
	})
}

func oneDishButton(rightPanel *fyne.Container, db *sql.DB) *widget.Button {
	return widget.NewButton("One Dish", func() {
		rightPanel.Refresh()
	})
}
