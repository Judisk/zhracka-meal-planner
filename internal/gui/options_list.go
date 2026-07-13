package gui

import (
	"database/sql"
	"math/rand/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func optionsList(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand, state *FilteredState) *fyne.Container {
	return container.NewVBox(
		allListButton(rightPanel, db, w, state),
		oneDayButton(rightPanel, db, w, rng),
		oneDishButton(rightPanel, db, w, rng),
	)
}

func allListButton(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, state *FilteredState) *widget.Button {
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
func oneDayButton(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *widget.Button {
	return widget.NewButton("One Day", func() {
		dayView := OneDayView(rightPanel, db, w, rng)
		rightPanel.Objects[0] = dayView
		rightPanel.Refresh()
	})
}

func oneDishButton(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *widget.Button {
	return widget.NewButton("One Dish", func() {
		dishView := OneDishView(rightPanel, db, w, rng)
		rightPanel.Objects[0] = dishView
		rightPanel.Refresh()
	})
}
