package gui

import (
	"database/sql"
	"foods/internal/dayone"
	"foods/internal/foodgenerator"
	"foods/internal/service"
	"math/rand/v2"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type dayState struct {
	n   *int
	day *dayone.Day
}

func OneDayView(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *fyne.Container {
	defN := 3
	defDay := dayone.Day{}
	dayState := dayState{
		n:   &defN,
		day: &defDay,
	}

	resultsPanel := container.NewVBox()

	buttonsContainer := container.NewVBox(
		generateButton(dayState, resultsPanel, db, w, rng),
		saveButton(dayState, db, w),
	)
	return container.NewBorder(
		nSelector(dayState, w),
		nil,
		nil,
		container.NewBorder(nil, buttonsContainer, nil, nil, nil),
		resultsPanel)
}

func nSelector(state dayState, w fyne.Window) *widget.Select {
	s := widget.NewSelect([]string{"1", "2", "3", "4", "5", "6"}, func(selected string) {
		num, err := strconv.Atoi(selected)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		*state.n = num
	})
	s.Selected = strconv.Itoa(*state.n)
	return s
}

func saveButton(state dayState, db *sql.DB, w fyne.Window) *widget.Button {
	return widget.NewButton("Save", func() {
		if len(state.day.Meals) == 0 {
			dialog.ShowInformation("Info", "Generate a day first", w)
			return
		}
		if err := service.ResetAndUpdateScore(db, *state.day); err != nil {
			dialog.ShowError(err, w)
			return
		}
	})
}

func generateButton(state dayState, resultsPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *widget.Button {
	return widget.NewButton(
		"Generate", func() {
			dayResult, err := service.GenerateDay(db, *state.n, rng)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			*state.day = dayResult
			resultsPanel.RemoveAll()
			resultsPanel.Add(manyMealsLabel(state))
			resultsPanel.Refresh()
		})
}

func oneMealLabel(grid *fyne.Container, dish foodgenerator.Dish) {

	grid.Add(widget.NewLabel(dish.Name))
	grid.Add(widget.NewLabel(dish.Grain.Name))
	grid.Add(widget.NewLabel(dish.Protein.Name))
	grid.Add(widget.NewLabel(dish.Vegetable.Name))

}

func manyMealsLabel(state dayState) *fyne.Container {

	grid := container.NewGridWithColumns(4)

	grid.Add(widget.NewLabel("Meal"))
	grid.Add(widget.NewLabel("Grain"))
	grid.Add(widget.NewLabel("Protein"))
	grid.Add(widget.NewLabel("Vegetable"))

	for _, elem := range state.day.Meals {
		oneMealLabel(grid, elem)
	}
	return grid
}
