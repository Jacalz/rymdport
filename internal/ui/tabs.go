// Package ui handles all logic related to the user interface.
package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/Jacalz/rymdport/v3/internal/transport"
)

// Create will set up and create the ui components.
func Create(app fyne.App, window fyne.Window) *container.AppTabs {
	client := transport.NewClient(app)

	tabs := &container.AppTabs{Items: []*container.TabItem{
		newSend(app, window, client).tabItem(),
		newRecv(app, window, client).tabItem(),
		newSettings(app, window, client).tabItem(),
		newAbout(app).tabItem(),
	}}

	// TODO: Tidy this up once we add mobile support.
	ctrlTab := &desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: fyne.KeyModifierControl}
	window.Canvas().AddShortcut(ctrlTab, func(shortcut fyne.Shortcut) {
		next := tabs.SelectedIndex() + 1
		if next >= len(tabs.Items) {
			next = 0
		}

		tabs.SelectIndex(next)
	})

	return tabs
}
