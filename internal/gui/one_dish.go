package gui

import (
	"database/sql"
	"foods/internal/dayone"
	"math/rand/v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func OneDishView(rightPanel *fyne.Container, db *sql.DB, w fyne.Window, rng *rand.Rand) *fyne.Container {
	n := 1
	defDay := dayone.Day{}
	defSaved := false
	dishState := mealsState{
		n:                  &n,
		day:                &defDay,
		savedStatus:        &defSaved,
		needName:           false,
		InfoTextSaveButton: "Generate a dish first",
	}
	resultsPanel := container.NewVBox()
	ondeDishButtonsContainer := container.NewVBox(
		generateButton(dishState, resultsPanel, db, w, rng),
		saveButton(dishState, db, w),
	)
	return container.NewBorder(
		nil, nil, nil, ondeDishButtonsContainer, resultsPanel,
	)
}
