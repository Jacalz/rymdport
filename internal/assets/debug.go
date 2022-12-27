//go:build !release && !debug

// Package assets contains bundled static resources.
package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed icon/icon-512.png
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
