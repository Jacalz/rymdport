//go:build !release && !debug

// Package assets contains bundled static resources.
package assets

import (
	_ "embed" // Needed for embedding the iconData on line 14.

	"fyne.io/fyne/v2"
)

//go:embed icons/icon-512.png
var iconData []byte

// SetIcon sets the icon because it is not bundled automatically.
func SetIcon(a fyne.App) {
	a.SetIcon(
		&fyne.StaticResource{
			StaticName:    "icon.png",
			StaticContent: iconData,
		},
	)
}
