package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func border(left fyne.CanvasObject, right fyne.CanvasObject) *container.Split {
	split := container.NewHSplit(left, right)
	split.SetOffset(0.2)
	return split
}
