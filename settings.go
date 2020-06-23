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
	PassPhraseComponentLength int

	// DownloadPath holds the download path used for saving recvieved files.
	DownloadPath string
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
	downloadPath.SetPlaceHolder("Downloads directory")
	downloadPath.OnChanged = func(input string) {
		switch input {
		case "":
			s.DownloadPath = userDownloadsFolder()
		default:
			// TODO: Make sure to only allow saving inside the home directory.
			s.DownloadPath = path.Clean(input)
		}
	}

	wormholeSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Download Path"), downloadPath)

	wormholeGroup := widget.NewGroup("Data Handling", wormholeSettingsContainer)

	settingsContent := widget.NewScrollContainer(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), interfaceGroup, layout.NewSpacer(), wormholeGroup))

	return widget.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settingsContent)
}
