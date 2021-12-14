package ui

import (
	"path/filepath"

	"fyne.io/fyne/v2"
	appearance "fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Jacalz/wormhole-gui/v2/internal/transport"
	"github.com/psanford/wormhole-william/wormhole"
)

type settings struct {
	downloadPathButton *widget.Button
	overwriteFiles     *widget.RadioGroup
	notificationRadio  *widget.RadioGroup

	componentSlider     *widget.Slider
	componentLabel      *widget.Label
	appID               *widget.Entry
	rendezvousURL       *widget.Entry
	transitRelayAddress *widget.Entry

	client      *transport.Client
	preferences fyne.Preferences
	window      fyne.Window
	app         fyne.App
}

func newSettings(a fyne.App, w fyne.Window, c *transport.Client) *settings {
	return &settings{app: a, window: w, client: c, preferences: a.Preferences()}
}

func (s *settings) onDownloadsPathChanged() {
	folder := dialog.NewFolderOpen(func(folder fyne.ListableURI, err error) {
		if err != nil {
			fyne.LogError("Error on selecting folder", err)
			dialog.ShowError(err, s.window)
			return
		} else if folder == nil {
			return
		}

		s.client.DownloadPath = folder.Path()
		s.preferences.SetString("DownloadPath", s.client.DownloadPath)
		s.downloadPathButton.SetText(folder.Name())
	}, s.window)

	folder.Resize(transport.WindowSizeToDialog(s.window.Canvas().Size()))
	folder.Show()
}

func (s *settings) onOverwriteFilesChanged(selected string) {
	s.client.OverwriteExisting = selected == "On"
	s.preferences.SetBool("OverwriteFiles", s.client.OverwriteExisting)
}

func (s *settings) onNotificationsChanged(selected string) {
	s.client.Notifications = selected == "On"
	s.preferences.SetBool("Notifications", s.client.Notifications)
}

func (s *settings) onComponentsChange(value float64) {
	s.client.PassPhraseComponentLength = int(value)
	s.preferences.SetInt("ComponentLength", s.client.PassPhraseComponentLength)
	s.componentLabel.SetText(string('0' + byte(value)))
}

func (s *settings) onAppIDChanged(appID string) {
	s.client.AppID = appID
	s.preferences.SetString("AppID", appID)
}

func (s *settings) onRendezvousURLChange(url string) {
	s.client.RendezvousURL = url
	s.preferences.SetString("RendezvousURL", url)
}

func (s *settings) onTransitAdressChange(address string) {
	s.client.TransitRelayAddress = address
	s.preferences.SetString("TransitRelayAddress", address)
}

// getPreferences is used to set the preferences on startup without saving at the same time.
func (s *settings) getPreferences() {
	s.client.DownloadPath = s.preferences.StringWithFallback("DownloadPath", transport.UserDownloadsFolder())
	s.downloadPathButton.Text = filepath.Base(s.client.DownloadPath)

	s.client.OverwriteExisting = s.preferences.Bool("OverwriteFiles")
	s.overwriteFiles.Selected = onOrOff(s.client.OverwriteExisting)

	s.client.Notifications = s.preferences.Bool("Notifications")
	s.notificationRadio.Selected = onOrOff(s.client.Notifications)

	s.client.PassPhraseComponentLength = s.preferences.IntWithFallback("ComponentLength", 2)
	s.componentSlider.Value = float64(s.client.PassPhraseComponentLength)
	s.componentLabel.Text = string('0' + byte(s.componentSlider.Value))

	s.client.AppID = s.preferences.String("AppID")
	s.appID.Text = s.client.AppID

	s.client.RendezvousURL = s.preferences.String("RendezvousURL")
	s.rendezvousURL.Text = s.client.RendezvousURL

	s.client.TransitRelayAddress = s.preferences.String("TransitRelayAddress")
	s.transitRelayAddress.Text = s.client.TransitRelayAddress
}

func (s *settings) buildUI() *container.Scroll {
	onOffOptions := []string{"On", "Off"}

	s.downloadPathButton = &widget.Button{Icon: theme.FolderOpenIcon(), OnTapped: s.onDownloadsPathChanged}

	s.overwriteFiles = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onOverwriteFilesChanged}

	s.notificationRadio = &widget.RadioGroup{Options: onOffOptions, Horizontal: true, Required: true, OnChanged: s.onNotificationsChanged}
	s.componentSlider, s.componentLabel = &widget.Slider{Min: 2.0, Max: 6.0, Step: 1, OnChanged: s.onComponentsChange}, &widget.Label{}

	s.appID = &widget.Entry{PlaceHolder: wormhole.WormholeCLIAppID, OnChanged: s.onAppIDChanged}

	s.rendezvousURL = &widget.Entry{PlaceHolder: wormhole.DefaultRendezvousURL, OnChanged: s.onRendezvousURLChange}

	s.transitRelayAddress = &widget.Entry{PlaceHolder: wormhole.DefaultTransitRelayAddress, OnChanged: s.onTransitAdressChange}

	s.getPreferences()

	interfaceContainer := appearance.NewSettings().LoadAppearanceScreen(s.window)

	dataContainer := container.NewGridWithColumns(2,
		newBoldLabel("Downloads Path"), s.downloadPathButton,
		newBoldLabel("Overwrite Files"), s.overwriteFiles,
		newBoldLabel("Notifications"), s.notificationRadio,
	)

	wormholeContainer := container.NewVBox(
		container.NewGridWithColumns(2, newBoldLabel("Passphrase Length"), container.NewBorder(nil, nil, nil, s.componentLabel, s.componentSlider)),
		&widget.Accordion{Items: []*widget.AccordionItem{
			{Title: "Advanced", Detail: container.NewGridWithColumns(2,
				newBoldLabel("AppID"), s.appID,
				newBoldLabel("Rendezvous URL"), s.rendezvousURL,
				newBoldLabel("Transit Relay Address"), s.transitRelayAddress,
			)},
		}},
	)

	return container.NewScroll(container.NewVBox(
		&widget.Card{Title: "User Interface", Content: interfaceContainer},
		&widget.Card{Title: "Data Handling", Content: dataContainer},
		&widget.Card{Title: "Wormhole Options", Content: wormholeContainer},
	))
}

func (s *settings) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Settings", Icon: theme.SettingsIcon(), Content: s.buildUI()}
}

func newBoldLabel(text string) *widget.Label {
	return &widget.Label{Text: text, TextStyle: fyne.TextStyle{Bold: true}}
}

func onOrOff(on bool) string {
	if on {
		return "On"
	}

	return "Off"
}
