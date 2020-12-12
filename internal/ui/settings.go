package ui

import (
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/bridge"
)

var (
	themes              = []string{"Adaptive (requires restart)", "Light", "Dark"}
	notificationOptions = []string{"On", "Off"}
)

// AppSettings cotains settings specific to the application
type AppSettings struct {
	// Theme holds the current theme
	Theme string
}

type settings struct {
	themeSelect *widget.Select

	downloadPathButton *widget.Button

	notificationRadio *widget.RadioGroup

	componentSlider *widget.Slider

	bridge      *bridge.Bridge
	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newSettings(a fyne.App, w fyne.Window, b *bridge.Bridge, as *AppSettings) *settings {
	return &settings{app: a, window: w, bridge: b, appSettings: as}
}

func (s *settings) onThemeChanged(selected string) {
	s.app.Preferences().SetString("Theme", checkTheme(selected, s.app))
}

func (s *settings) onDownloadsPathChanged() {
	dialog.ShowFolderOpen(func(folder fyne.ListableURI, err error) {
		if err != nil {
			fyne.LogError("Error on selecting folder", err)
			dialog.ShowError(err, s.window)
			return
		} else if folder == nil {
			return
		}

		s.app.Preferences().SetString("DownloadPath", folder.String()[7:])
		s.bridge.DownloadPath = folder.String()[7:]
		s.downloadPathButton.SetText(folder.Name())
	}, s.window)
}

func (s *settings) onNotificationsChanged(selected string) {
	if selected == "On" {
		s.bridge.Notifications = true
	} else {
		s.bridge.Notifications = false
	}

	s.app.Preferences().SetString("Notifications", selected)
}

func (s *settings) onComponentsChange(value float64) {
	s.bridge.PassPhraseComponentLength = int(value)
	s.app.Preferences().SetFloat("ComponentLength", value)
}

func (s *settings) buildUI() *container.Scroll {
	s.themeSelect = &widget.Select{Options: themes, OnChanged: s.onThemeChanged, Selected: s.appSettings.Theme}

	s.bridge.DownloadPath = s.app.Preferences().StringWithFallback("DownloadPath", bridge.UserDownloadsFolder())
	s.downloadPathButton = &widget.Button{Icon: theme.FolderOpenIcon(), OnTapped: s.onDownloadsPathChanged, Text: filepath.Base(s.bridge.DownloadPath)}

	s.notificationRadio = &widget.RadioGroup{Options: notificationOptions, Horizontal: true, Required: true, OnChanged: s.onNotificationsChanged}
	s.notificationRadio.SetSelected(s.app.Preferences().StringWithFallback("Notifications", notificationOptions[1]))

	s.componentSlider = &widget.Slider{Min: 2.0, Max: 6.0, OnChanged: s.onComponentsChange}
	s.componentSlider.SetValue(s.app.Preferences().FloatWithFallback("ComponentLength", 2))

	interfaceContainer := container.NewGridWithColumns(2,
		newSettingLabel("Application Theme"), s.themeSelect,
	)
	interfaceCard := widget.NewCard("User Interface", "Settings to manage the application appearance.", interfaceContainer)

	dataContainer := container.NewGridWithColumns(2,
		newSettingLabel("Downloads Path"), s.downloadPathButton,
		newSettingLabel("Notifications"), s.notificationRadio,
	)
	dataCard := widget.NewCard("Data Handling", "Settings for handling of data.", dataContainer)

	wormholeContainer := container.NewGridWithColumns(2,
		newSettingLabel("Passphrase Length"), s.componentSlider,
	)
	wormholeCard := widget.NewCard("Wormhole Options", "Settings for configuring wormhole.", wormholeContainer)

	return container.NewScroll(container.NewVBox(
		interfaceCard,
		dataCard,
		wormholeCard,
	))
}

func (s *settings) tabItem() *container.TabItem {
	return container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), s.buildUI())
}

func checkTheme(themec string, a fyne.App) string {
	switch themec {
	case "Light":
		a.Settings().SetTheme(theme.LightTheme())
	case "Dark":
		a.Settings().SetTheme(theme.DarkTheme())
	}

	return themec
}

func newSettingLabel(text string) *widget.Label {
	return &widget.Label{TextStyle: fyne.TextStyle{Bold: true}, Text: text}
}
