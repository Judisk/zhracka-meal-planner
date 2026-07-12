package gui

import (
	"database/sql"
	"foods/internal/dayone"
	"foods/internal/foodgenerator"
	"foods/internal/service"
	"math/rand/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type mealsState struct {
	n                  *int
	day                *dayone.Day
	savedStatus        *bool
	needName           bool
	InfoTextSaveButton string
}

func saveButton(state mealsState, db *sql.DB, w fyne.Window) *widget.Button {
	return widget.NewButton("Save", func() {
		if len(state.day.Meals) == 0 {
			dialog.ShowInformation("Info", state.InfoTextSaveButton, w)
			return
		}
		if *state.savedStatus {
			dialog.ShowInformation("Info", "Already saved", w)
			return
		}
		if err := service.ResetAndUpdateScore(db, *state.day); err != nil {
			dialog.ShowError(err, w)
			return
		}
		*state.savedStatus = true
	})
}

func generateButton(state mealsState, resultsPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *widget.Button {
	return widget.NewButton(
		"Generate", func() {
			dayResult, err := service.GenerateDay(db, *state.n, rng)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			*state.day = dayResult
			*state.savedStatus = false
			resultsPanel.RemoveAll()
			resultsPanel.Add(manyMealsLabel(state))
			resultsPanel.Refresh()
		})
}
func oneMealLabel(grid *fyne.Container, dish foodgenerator.Dish, needName bool) {
	if needName {
		grid.Add(widget.NewLabel(dish.Name))
	}
	grid.Add(widget.NewLabel(dish.Grain.Name))
	grid.Add(widget.NewLabel(dish.Protein.Name))
	grid.Add(widget.NewLabel(dish.Vegetable.Name))
}

func manyMealsLabel(state mealsState) *fyne.Container {
	var columnCount int
	if state.needName {
		columnCount = 4
	} else {
		columnCount = 3
	}

	grid := container.NewGridWithColumns(columnCount)

	if state.needName {
		grid.Add(widget.NewLabel("Meal"))
	}
	grid.Add(widget.NewLabel("Grain"))
	grid.Add(widget.NewLabel("Protein"))
	grid.Add(widget.NewLabel("Vegetable"))

	for _, elem := range state.day.Meals {
		oneMealLabel(grid, elem, state.needName)
	}
	return grid
}
