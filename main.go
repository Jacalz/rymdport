package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/Jacalz/wormhole-gui/assets"
)

func main() {
	a := app.NewWithID("com.github.jacalz.wormhole-gui")
	a.SetIcon(assets.AppIcon)
	w := a.NewWindow("wormhole-gui")

	switch a.Preferences().StringWithFallback("Theme", "Light") {
	case "Light":
		a.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		a.Settings().SetTheme(theme.DarkTheme())
	}

	s := settings{ComponentLength: 2, DownloadPath: userDownloadsFolder()}

	w.SetContent(widget.NewTabContainer(s.sendTab(a, w), s.recieveTab(w), s.settingsTab(a), aboutTab()))
	w.Resize(fyne.NewSize(600, 400))
	w.ShowAndRun()
}
