package gui

import (
	"database/sql"
	"foods/internal/dayone"
	"foods/internal/service"
	"math/rand/v2"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func OneDayView(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *fyne.Container {
	defN := 3
	defDay := dayone.Day{}
	defSaved := false
	dayState := mealsState{
		n:                  &defN,
		day:                &defDay,
		savedStatus:        &defSaved,
		needName:           true,
		InfoTextSaveButton: "Generate a day first",
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

func nSelector(state mealsState, w fyne.Window) *widget.Select {
	s := widget.NewSelect(service.ConvertMinToMaxInString(service.MinMeals, service.MaxMeals), func(selected string) {
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
