package bridge

import (
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"github.com/mholt/archiver/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

// Bridge holds settings and other meathods specific to the bridge using wormhole-willam.
type Bridge struct {
	wormhole.Client

	// Save a reference to the window to avoid creating a new one when sending and receiving text
	window fyne.Window

	// Save a reference to the zip handler to avoid creating a new one each time when unzipping folders
	zip *archiver.Zip

	// Notification holds the settings value for if we have notifications enabled or not.
	Notifications bool

	// DownloadPath holds the download path used for saving recvieved files.
	DownloadPath string
}

// NewBridge returns a new bridge that is configured and ready
func NewBridge() *Bridge {
	b := &Bridge{window: fyne.CurrentApp().NewWindow(""), zip: &archiver.Zip{MkdirAll: true}}
	b.window.Resize(fyne.NewSize(400, 300))
	return b
}

// UserDownloadsFolder returns the downloads folder corresponding to the current user.
func UserDownloadsFolder() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not get home dir", err)
	}

	return filepath.Join(dir, "Downloads")
}
