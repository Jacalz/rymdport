package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/Jacalz/rymdport/v3/internal/assets"
	"github.com/Jacalz/rymdport/v3/internal/ui"
)

func main() {
	a := app.NewWithID("io.github.jacalz.rymdport")
	assets.SetIcon(a)
	w := a.NewWindow("Rymdport")

	w.SetContent(ui.Create(a, w))
	w.Resize(fyne.NewSize(600, 400))
	w.SetMaster()
	w.ShowAndRun()
}
