package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func main() {
	a := app.NewWithID("com.github.jacalz.wormhole-gui")
	w := a.NewWindow("wormhole-gui")

	switch a.Preferences().StringWithFallback("Theme", "Light") {
	case "Dark":
		a.Settings().SetTheme(theme.DarkTheme())
	case "Light":
		a.Settings().SetTheme(theme.LightTheme())
	}

	w.SetContent(widget.NewTabContainer(sendTab(), recieveTab(), settingsTab()))
	w.ShowAndRun()
}
