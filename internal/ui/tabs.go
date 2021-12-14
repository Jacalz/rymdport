package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Jacalz/wormhole-gui/v2/internal/transport"
)

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	bridge := transport.NewClient(app)

	return &container.AppTabs{Items: []*container.TabItem{
		newSend(app, window, bridge).tabItem(),
		newRecv(app, window, bridge).tabItem(),
		newSettings(app, window, bridge).tabItem(),
		newAbout().tabItem(),
	}}
}
