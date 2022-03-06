package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Jacalz/rymdport/v3/internal/transport"
)

// Create will stitch together all ui components
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	client := transport.NewClient(app)

	return &container.AppTabs{Items: []*container.TabItem{
		newSend(app, window, client).tabItem(),
		newRecv(app, window, client).tabItem(),
		newSettings(app, window, client).tabItem(),
		newAbout().tabItem(),
	}}
}
