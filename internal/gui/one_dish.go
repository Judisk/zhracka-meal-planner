package gui

import (
	"database/sql"
	"foods/internal/dayone"
	"foods/internal/service"
	"math/rand/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func OneDishView(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *fyne.Container {
	result := &dayone.Day{}
	resultsPanel := container.NewVBox()
	ondeDishButtonsContainer := container.NewVBox(
		generateOneDishButton(result, resultsPanel, db, w, rng),
		oneDishSaveButton(result, db, w),
	)
	return container.NewBorder(
		nil, nil, nil, ondeDishButtonsContainer, resultsPanel,
	)
}

func generateOneDishButton(result *dayone.Day, resultsPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *widget.Button {
	return widget.NewButton("Generate", func() {
		dayResult, err := service.GenerateDay(db, 1, rng)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		*result = dayResult
		resultsPanel.RemoveAll()
		resultsPanel.Add(oneDishLabel(*result))
		resultsPanel.Refresh()
	})
}

func oneDishLabel(state dayone.Day) *fyne.Container {
	dish := state.Meals[0]
	grid := container.NewGridWithColumns(3)

	grid.Add(widget.NewLabel("Grain"))
	grid.Add(widget.NewLabel("Protein"))
	grid.Add(widget.NewLabel("Vegetable"))

	grid.Add(widget.NewLabel(dish.Grain.Name))
	grid.Add(widget.NewLabel(dish.Protein.Name))
	grid.Add(widget.NewLabel(dish.Vegetable.Name))
	return grid
}

func oneDishSaveButton(state *dayone.Day, db *sql.DB, w fyne.Window) *widget.Button {
	return widget.NewButton("Save", func() {
		if len(state.Meals) == 0 {
			dialog.ShowInformation("Info", "Generate a day first", w)
			return
		}
		if err := service.ResetAndUpdateScore(db, *state); err != nil {
			dialog.ShowError(err, w)
			return
		}
	})
}
