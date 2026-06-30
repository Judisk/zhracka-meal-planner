package gui

import (
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func Run(db *sql.DB) {
	a := app.New()
	w := a.NewWindow("Food Planner")
	w.Resize(fyne.NewSize(sumOfSize()+300, 500))

	productTableContainer, err := productsTable(db)
	if err != nil {
		// TODO: handle the error in the GUI
	}
	rightPanel := container.NewStack(productTableContainer)
	optionListContainer := optionsList(rightPanel, db)
	content := border(optionListContainer, rightPanel)
	w.SetContent(content)
	w.ShowAndRun()
}
