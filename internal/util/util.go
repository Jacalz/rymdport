// Package util contains various small helper functions.
package util

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
)

//lint:ignore ST1005 The error is not printed to a terminal as usual but displayed to the user in the ui.
var errInvalidCode = errors.New("Invalid code. Must begin with a number followed by groups of letters, separated with \"-\".")

// CodeValidator provides a validator for wormhole codes.
func CodeValidator(input string) error {
	if input == "" {
		return nil // We don't want empty entry to report an error.
	}

	next := strings.IndexByte(input, '-')
	if next <= 0 || next == len(input)-1 {
		return errInvalidCode
	}

	invalidNameplate := strings.IndexFunc(input[:next], runeIsNotNumerical) != -1
	if invalidNameplate {
		return errInvalidCode
	}

	invalidChars := strings.IndexFunc(input[next+1:], runeIsInvalid) != -1
	if invalidChars {
		return errInvalidCode
	}

	for next != -1 {
		input = input[next+1:]

		next = strings.IndexByte(input, '-')
		if next == 0 || next == len(input)-1 {
			return errInvalidCode
		}
	}

	return nil
}

func runeIsInvalid(r rune) bool {
	return (r < '0' || r > '9') && (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && r != '-'
}

func runeIsNotNumerical(r rune) bool {
	return r < '0' || r > '9'
}

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
