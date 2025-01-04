package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

func createRecvPage(_ fyne.App, _ fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("Receive something")
}
