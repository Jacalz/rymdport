//go:build !release

package updater

import "fyne.io/fyne/v2"

// Enable is turned off when not compiling release binaries.
func Enable(_ fyne.App, _ fyne.Window) {}
