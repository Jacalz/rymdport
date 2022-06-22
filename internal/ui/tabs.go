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

	// Set up support for cycling through the tabs.
	ctrlTab := &desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: fyne.KeyModifierControl}
	window.Canvas().AddShortcut(ctrlTab, func(_ fyne.Shortcut) {
		next := tabs.SelectedIndex() + 1
		if next >= len(tabs.Items) {
			next = 0
		}

		tabs.SelectIndex(next)
	})

	// Set up support for Alt + [1:4] for switching to a specific tab.
	tabnumber := func(shortcut fyne.Shortcut) {
		name := shortcut.ShortcutName()
		tabs.SelectIndex(int(name[len(name)-1]-'0') - 1)
	}

	alt1 := &desktop.CustomShortcut{KeyName: fyne.Key1, Modifier: fyne.KeyModifierAlt}
	window.Canvas().AddShortcut(alt1, tabnumber)
	alt2 := &desktop.CustomShortcut{KeyName: fyne.Key2, Modifier: fyne.KeyModifierAlt}
	window.Canvas().AddShortcut(alt2, tabnumber)
	alt3 := &desktop.CustomShortcut{KeyName: fyne.Key3, Modifier: fyne.KeyModifierAlt}
	window.Canvas().AddShortcut(alt3, tabnumber)
	alt4 := &desktop.CustomShortcut{KeyName: fyne.Key4, Modifier: fyne.KeyModifierAlt}
	window.Canvas().AddShortcut(alt4, tabnumber)

	return tabs
}
