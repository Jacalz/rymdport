package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"github.com/Jacalz/wormhole-gui/internal/bridge"
)

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	bridge := bridge.NewBridge() // To make sure that it is configured correctly
	appSettings := &AppSettings{}
	appSettings.Theme = checkTheme(app.Preferences().StringWithFallback("Theme", "Adaptive (requires restart)"), app)

	return container.NewAppTabs(
		newSend(app, window, bridge, appSettings).tabItem(),
		newRecv(app, window, bridge, appSettings).tabItem(),
		newSettings(app, window, bridge, appSettings).tabItem(),
		newAbout().tabItem(),
	)
}
