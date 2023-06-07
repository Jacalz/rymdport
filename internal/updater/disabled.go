//go:build !release

package updater

import "fyne.io/fyne/v2"

// Enabled specifies if the binary is built with update checking enabled.
const Enabled = false

// Enable is turned off when not compiling release binaries.
func Enable(_ fyne.App, _ fyne.Window) {}
