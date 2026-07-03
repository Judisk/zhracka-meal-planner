package gui

import (
	"database/sql"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Run(db *sql.DB) {
	a := app.New()
	w := a.NewWindow("Food Planner")
	w.Resize(fyne.NewSize(sumOfSize()+300, 500))
	rightPanel := container.NewStack(widget.NewLabel("Loading..."))

	productTableContainer, err := productsTable(db, w, rightPanel)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	rightPanel.Objects[0] = productTableContainer
	rightPanel.Refresh()

	optionListContainer := optionsList(rightPanel, db, w)
	content := border(optionListContainer, rightPanel)
	w.SetContent(content)
	w.ShowAndRun()
}
