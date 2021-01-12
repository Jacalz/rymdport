package ui

import (
	"path/filepath"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/transport"
)

var (
	themes       = []string{"Adaptive (requires restart)", "Light", "Dark"}
	onOffOptions = []string{"On", "Off"}
)

// AppSettings contains settings specific to the application
type AppSettings struct {
	// Theme holds the current theme
	Theme string
}

type settings struct {
	themeSelect *widget.Select

	downloadPathButton *widget.Button
	overwriteFiles     *widget.RadioGroup
	notificationRadio  *widget.RadioGroup

	componentSlider *widget.Slider
	componentLabel  *widget.Label

	client      *transport.Client
	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newSettings(a fyne.App, w fyne.Window, c *transport.Client, as *AppSettings) *settings {
	return &settings{app: a, window: w, client: c, appSettings: as}
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
		s.client.DownloadPath = folder.String()[7:]
		s.downloadPathButton.SetText(folder.Name())
	}, s.window)
}

func (s *settings) onOverwriteFilesChanged(selected string) {
	s.client.Zip.OverwriteExisting = selected == "On"
	s.app.Preferences().SetString("OverwriteFiles", selected)
}

func (s *settings) onNotificationsChanged(selected string) {
	s.client.Notifications = selected == "On"
	s.app.Preferences().SetString("Notifications", selected)
}

func (s *settings) onComponentsChange(value float64) {
	s.client.PassPhraseComponentLength = int(value)
	s.app.Preferences().SetFloat("ComponentLength", value)
	s.componentLabel.SetText(strconv.Itoa(int(value)))
}

func (s *settings) buildUI() *container.Scroll {
	s.themeSelect = &widget.Select{Options: themes, OnChanged: s.onThemeChanged, Selected: s.appSettings.Theme}

	s.client.DownloadPath = s.app.Preferences().StringWithFallback("DownloadPath", transport.UserDownloadsFolder())
	s.downloadPathButton = &widget.Button{Icon: theme.FolderOpenIcon(), OnTapped: s.onDownloadsPathChanged, Text: filepath.Base(s.client.DownloadPath)}

	s.overwriteFiles = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onOverwriteFilesChanged}
	s.overwriteFiles.SetSelected(s.app.Preferences().StringWithFallback("OverwriteFiles", "Off"))

	s.notificationRadio = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onNotificationsChanged}
	s.notificationRadio.SetSelected(s.app.Preferences().StringWithFallback("Notifications", onOffOptions[1]))

	s.componentLabel = &widget.Label{}
	s.componentSlider = &widget.Slider{Min: 2.0, Max: 6.0, Step: 1, OnChanged: s.onComponentsChange}
	s.componentSlider.SetValue(s.app.Preferences().FloatWithFallback("ComponentLength", 2))

	interfaceContainer := container.NewGridWithColumns(2,
		newSettingLabel("Application Theme"), s.themeSelect,
	)
	interfaceCard := widget.NewCard("User Interface", "Settings to manage the application appearance.", interfaceContainer)

	dataContainer := container.NewGridWithColumns(2,
		newSettingLabel("Downloads Path"), s.downloadPathButton,
		newSettingLabel("Overwrite Files"), s.overwriteFiles,
		newSettingLabel("Notifications"), s.notificationRadio,
	)
	dataCard := widget.NewCard("Data Handling", "Settings for handling of data.", dataContainer)

	wormholeContainer := container.NewGridWithColumns(2,
		newSettingLabel("Passphrase Length"), container.NewBorder(nil, nil, nil, s.componentLabel, s.componentSlider),
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
