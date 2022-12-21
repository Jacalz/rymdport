package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/Jacalz/rymdport/v3/internal/assets"
	"github.com/Jacalz/rymdport/v3/internal/ui"
	"github.com/Jacalz/rymdport/v3/internal/updater"
)

func main() {
	a := app.NewWithID("io.github.jacalz.rymdport")
	a.SetIcon(assets.Icon)
	w := a.NewWindow("Rymdport")

	w.SetContent(ui.Create(a, w))
	w.Resize(fyne.NewSize(700, 400))
	w.SetMaster()

	updater.Enable(a, w)
	w.ShowAndRun()
}
