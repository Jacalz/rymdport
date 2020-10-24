package main

import (
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/bridge"
)

func checkTheme(themec string, a fyne.App) string {
	switch themec {
	case "Light":
		a.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		a.Settings().SetTheme(theme.DarkTheme())
	}

	return themec
}

func (ad *appData) settingsTab() *container.TabItem {
	themeSwitcher := widget.NewSelect([]string{"Adaptive (requires restart)", "Light", "Dark"}, func(selected string) {
		ad.App.Preferences().SetString("Theme", checkTheme(selected, ad.App))
	})
	themeSwitcher.SetSelected(ad.App.Preferences().StringWithFallback("Theme", "Adaptive (requires restart)"))

	interfaceSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Application Theme"), themeSwitcher)

	interfaceGroup := widget.NewGroup("User Interface", interfaceSettingsContainer)

	downloadPathButton := &widget.Button{Text: "Downloads", Icon: theme.FolderOpenIcon()}
	ad.Bridge.DownloadPath = ad.App.Preferences().StringWithFallback("DownloadPath", bridge.UserDownloadsFolder())
	downloadPathButton.SetText(filepath.Base(ad.Bridge.DownloadPath))
	downloadPathButton.OnTapped = func() {
		dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				fyne.LogError("Error on selecting folder", err)
				dialog.ShowError(err, ad.Window)
				return
			} else if folder == nil {
				return
			}

			ad.App.Preferences().SetString("DownloadPath", folder.String())
			ad.Bridge.DownloadPath = folder.String()[7:]
			downloadPathButton.SetText(folder.Name())
		}, ad.Window)
	}

	notification := &widget.Radio{Options: []string{"On", "Off"}, Horizontal: true, Required: true, OnChanged: func(selected string) {
		if selected == "On" {
			ad.Notifications = true
		} else {
			ad.Notifications = false
		}

		ad.App.Preferences().SetString("Notifications", selected)
	}}
	notification.SetSelected(ad.App.Preferences().StringWithFallback("Notifications", "Off"))

	dataSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Download Path"), downloadPathButton, widget.NewLabel("Notifications"), notification)
	dataGroup := widget.NewGroup("Data Handling", dataSettingsContainer)

	slider := widget.NewSlider(2.0, 6.0)
	slider.OnChanged = func(value float64) {
		ad.Bridge.ComponentLength = int(value)
		ad.App.Preferences().SetFloat("ComponentLength", float64(ad.Bridge.ComponentLength))
	}
	slider.SetValue(ad.App.Preferences().FloatWithFallback("ComponentLength", 2))

	wormholeSettingsContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widget.NewLabel("Passphrase Length"), slider)
	wormholeGroup := widget.NewGroup("Wormhole Parameters", wormholeSettingsContainer)

	settingsContent := container.NewScroll(fyne.NewContainerWithLayout(layout.NewVBoxLayout(), interfaceGroup, layout.NewSpacer(), dataGroup, layout.NewSpacer(), wormholeGroup))

	return widget.NewTabItemWithIcon("Settings", theme.SettingsIcon(), settingsContent)
}
