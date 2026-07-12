package gui

import (
	"database/sql"
	"math/rand/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Run(db *sql.DB, rng *rand.Rand) {
	a := app.New()
	w := a.NewWindow("Food Planner")
	w.Resize(fyne.NewSize(sumOfSize()+300, 500))
	rightPanel := container.NewStack(widget.NewLabel("Loading..."))

	state := FilteredState{
		CategoryState:      nil,
		BannedState:        nil,
		PreferencesState:   nil,
		CategorySelected:   "Category",
		BannedSelected:     "All",
		PreferenceSelected: "Preference",
	}

	productTableContainer, err := tableContainer(db, w, rightPanel, state)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	rightPanel.Objects[0] = productTableContainer
	rightPanel.Refresh()

	optionListContainer := optionsList(rightPanel, db, w, rng, state)
	content := border(optionListContainer, rightPanel)
	w.SetContent(content)
	w.ShowAndRun()
}
