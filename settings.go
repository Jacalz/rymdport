package main

import (
	"os"
	"path"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func userDownloadsFolder() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not get home dir", err)
	}

	return path.Join(dir, "Downloads")
}

type settings struct {
	// PassPhraseComponentLength is the number of words to use when generating a passprase.
	ComponentLength int

	// DownloadPath holds the download path used for saving recvieved files.
	DownloadPath string

	// Notification holds the settings value for if we have notifications enabled or not.
	Notifications bool
}

func (s *settings) settingsTab(a fyne.App) *widget.TabItem {
	themeSwitcher := widget.NewSelect([]string{"Light", "Dark"}, func(selected string) {
		switch selected {
		case "Light":
			a.Settings().SetTheme(theme.LightTheme())
		case "Dark":
			a.Settings().SetTheme(theme.DarkTheme())
		}

		// Set the theme to the selected one and save it using the preferences api in fyne.
		a.Preferences().SetString("Theme", selected)
	})
	themeSwitcher.SetSelected(a.Preferences().StringWithFallback("Theme", "Light"))

	interfaceSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Application Theme"), themeSwitcher)

	interfaceGroup := widget.NewGroup("User Interface", interfaceSettingsContainer)

	downloadPath := widget.NewEntry()
	downloadPath.SetPlaceHolder("Downloads")
	downloadPath.OnChanged = func(input string) {
		switch input {
		case "":
			s.DownloadPath = userDownloadsFolder()
		default:
			// TODO: Make sure to only allow saving inside the home directory.
			s.DownloadPath = path.Clean(input)
		}
	}

	notification := widget.NewRadio([]string{"On", "Off"}, func(selected string) {
		if selected == "On" {
			s.Notifications = true
		} else {
			s.Notifications = false
		}

		a.Preferences().SetString("Notifications", selected)
	})

	notification.SetSelected(a.Preferences().StringWithFallback("Notifications", "Off"))
	notification.Horizontal = true

	dataSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Download Path"), downloadPath, widget.NewLabel("Notifications"), notification)
	dataGroup := widget.NewGroup("Data Handling", dataSettingsContainer)

	slider := widget.NewSlider(2.0, 6.0)
	slider.OnChanged = func(value float64) {
		s.ComponentLength = int(value)
	}
	slider.Refresh()

	wormholeSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Passphrase Length"), slider)
	wormholeGroup := widget.NewGroup("Wormhole Parameters", wormholeSettingsContainer)

	settingsContent := widget.NewScrollContainer(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), interfaceGroup, layout.NewSpacer(), dataGroup, layout.NewSpacer(), wormholeGroup))

	return widget.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settingsContent)
}
