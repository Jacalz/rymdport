package util

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
)

// CodeValidator provides a validator for wormhole codes.
var CodeValidator = validation.NewRegexp(`(^\d+(-(\w|\d)+)+$)|(^$)`, "The code string is invalid.")

// UserDownloadsFolder returns the downloads folder corresponding to the current user.
func UserDownloadsFolder() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not get home dir", err)
	}

	return filepath.Join(dir, "Downloads")
}

// WindowSizeToDialog scales the window size to a suitable dialog size.
func WindowSizeToDialog(s fyne.Size) fyne.Size {
	return fyne.NewSize(s.Width*0.8, s.Height*0.8)
}
