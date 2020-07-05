package main

import (
	"path"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/bridge"
)

func (ad *appData) settingsTab() *widget.TabItem {
	themeSwitcher := widget.NewSelect([]string{"Light", "Dark"}, func(selected string) {
		switch selected {
		case "Light":
			ad.App.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			ad.App.Settings().SetTheme(theme.DarkTheme())
		}

		// Set the theme to the selected one and save it using the preferences api in fyne.
		ad.App.Preferences().SetString("Theme", selected)
	})
	themeSwitcher.SetSelected(ad.App.Preferences().StringWithFallback("Theme", "Light"))

	interfaceSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Application Theme"), themeSwitcher)

	interfaceGroup := widget.NewGroup("User Interface", interfaceSettingsContainer)

	downloadPath := widget.NewEntry()
	downloadPath.SetPlaceHolder("Downloads")
	downloadPath.OnChanged = func(input string) {
		switch input {
		case "":
			ad.Bridge.DownloadPath = bridge.UserDownloadsFolder()
		default:
			// TODO: Make sure to only allow saving inside the home directory.
			ad.Bridge.DownloadPath = path.Clean(input)
		}
	}

	notification := widget.NewRadio([]string{"On", "Off"}, func(selected string) {
		if selected == "On" {
			ad.Notifications = true
		} else {
			ad.Notifications = false
		}

		ad.App.Preferences().SetString("Notifications", selected)
	})

	notification.SetSelected(ad.App.Preferences().StringWithFallback("Notifications", "Off"))
	notification.Horizontal = true

	dataSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Download Path"), downloadPath, widget.NewLabel("Notifications"), notification)
	dataGroup := widget.NewGroup("Data Handling", dataSettingsContainer)

	slider := widget.NewSlider(2.0, 6.0)
	slider.OnChanged = func(value float64) {
		ad.Bridge.ComponentLength = int(value)
	}

	wormholeSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Passphrase Length"), slider)
	wormholeGroup := widget.NewGroup("Wormhole Parameters", wormholeSettingsContainer)

	settingsContent := widget.NewScrollContainer(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), interfaceGroup, layout.NewSpacer(), dataGroup, layout.NewSpacer(), wormholeGroup))

	return widget.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settingsContent)
}
