package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"

	"github.com/Jacalz/wormhole-gui/assets"
	"github.com/Jacalz/wormhole-gui/bridge"
)

type appData struct {
	// BridgeSettings holds the settings specific to the bridge implementation.
	Bridge *bridge.Bridge

	// Notification holds the settings value for if we have notifications enabled or not.
	Notifications bool

	// App and window are variables from fyne.
	App    fyne.App
	Window fyne.Window
}

func main() {
	ad := appData{Bridge: bridge.NewBridge()}

	ad.App = app.NewWithID("com.github.jacalz.wormhole-gui")
	ad.App.SetIcon(assets.AppIcon)
	ad.Window = ad.App.NewWindow("wormhole-gui")

	checkTheme(ad.App.Preferences().StringWithFallback("Theme", "Adaptive (requires restart)"), ad.App)

	ad.Window.SetContent(widget.NewTabContainer(ad.sendTab(), ad.recieveTab(), ad.settingsTab(), aboutTab()))
	ad.Window.Resize(fyne.NewSize(600, 400))
	ad.Window.ShowAndRun()
}
