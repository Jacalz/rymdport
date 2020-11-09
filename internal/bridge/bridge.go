package bridge

import (
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"github.com/psanford/wormhole-william/wormhole"
)

// Bridge holds settings and other meathods specific to the bridge using wormhole-willam.
type Bridge struct {
	wormhole.Client

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
