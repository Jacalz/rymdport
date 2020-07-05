package bridge

import (
	"os"
	"path/filepath"

	"fyne.io/fyne"
)

// Bridge holds settings and other meathods specific to the bridge using wormhole-willam.
type Bridge struct {
	// PassPhraseComponentLength is the number of words to use when generating a passprase.
	ComponentLength int

	// DownloadPath holds the download path used for saving recvieved files.
	DownloadPath string
}

// UserDownloadsFolder returns the downloads folder corresponding to the current user.
func UserDownloadsFolder() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not get home dir", err)
	}

	return filepath.Join(dir, "Downloads")
}

// NewBridge creates a new bridge with the default values.
func NewBridge() *Bridge {
	return &Bridge{ComponentLength: 2, DownloadPath: UserDownloadsFolder()}
}
