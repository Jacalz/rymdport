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
	client := transport.NewClient(app, window)

	tabs := &container.AppTabs{
		Items: []*container.TabItem{
			newSendTab(window, client),
			newRecvTab(window, client),
			newSettingsTab(app, window, client),
			newAboutTab(app),
		},
	}

	canvas := window.Canvas()

	// Set up support for switching between the tabs.
	ctrlTab := &desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: fyne.KeyModifierControl}
	canvas.AddShortcut(ctrlTab, func(_ fyne.Shortcut) {
		next := tabs.SelectedIndex() + 1
		if next >= len(tabs.Items) {
			next = 0
		}

		tabs.SelectIndex(next)
	})

	ctrlShiftTab := &desktop.CustomShortcut{KeyName: fyne.KeyTab, Modifier: fyne.KeyModifierControl | fyne.KeyModifierShift}
	canvas.AddShortcut(ctrlShiftTab, func(_ fyne.Shortcut) {
		next := tabs.SelectedIndex() - 1
		if next < 0 {
			next += len(tabs.Items)
		}

		tabs.SelectIndex(next)
	})

	// Set up support for Alt + [1:4] for switching to a specific tab.
	alt1 := &desktop.CustomShortcut{KeyName: fyne.Key1, Modifier: fyne.KeyModifierAlt}
	canvas.AddShortcut(alt1, func(_ fyne.Shortcut) { tabs.SelectIndex(0) })
	alt2 := &desktop.CustomShortcut{KeyName: fyne.Key2, Modifier: fyne.KeyModifierAlt}
	canvas.AddShortcut(alt2, func(_ fyne.Shortcut) { tabs.SelectIndex(1) })
	alt3 := &desktop.CustomShortcut{KeyName: fyne.Key3, Modifier: fyne.KeyModifierAlt}
	canvas.AddShortcut(alt3, func(_ fyne.Shortcut) { tabs.SelectIndex(2) })
	alt4 := &desktop.CustomShortcut{KeyName: fyne.Key4, Modifier: fyne.KeyModifierAlt}
	canvas.AddShortcut(alt4, func(_ fyne.Shortcut) { tabs.SelectIndex(3) })

	return tabs
}
