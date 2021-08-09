package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"github.com/Jacalz/wormhole-gui/v2/internal/assets"
	"github.com/Jacalz/wormhole-gui/v2/internal/ui"
)

func main() {
	a := app.NewWithID("com.github.jacalz.wormhole-gui")
	a.SetIcon(assets.AppIcon)
	w := a.NewWindow("wormhole-gui")

	w.SetContent(ui.Create(a, w))
	w.Resize(fyne.NewSize(700, 400))
	w.SetMaster()
	w.ShowAndRun()
}
