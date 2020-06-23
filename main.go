package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.NewWithID("com.github.jacalz.wormhole-gui")
	w := a.NewWindow("wormhole-gui")

	switch a.Preferences().StringWithFallback("Theme", "Light") {
	case "Light":
		a.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		a.Settings().SetTheme(theme.DarkTheme())
	}

	w.SetContent(widget.NewTabContainer(sendTab(w), recieveTab(), settingsTab(a)))
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
