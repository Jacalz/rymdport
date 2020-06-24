package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func boldLabel(text string) *widget.Label {
	return widget.NewLabelWithStyle(text, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
}

func main() {
	a := app.NewWithID("com.github.jacalz.wormhole-gui")
	w := a.NewWindow("wormhole-gui")

	switch a.Preferences().StringWithFallback("Theme", "Light") {
	case "Light":
		a.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		a.Settings().SetTheme(theme.DarkTheme())
	}

	s := settings{ComponentLength: 2, DownloadPath: userDownloadsFolder()}

	w.SetContent(widget.NewTabContainer(s.sendTab(w), s.recieveTab(w), s.settingsTab(a)))
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
