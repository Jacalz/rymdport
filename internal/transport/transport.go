// Package transport handles sending and receiving using wormhole-william
package transport

import (
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"github.com/mholt/archiver/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

// Client defines the client for handling sending and receiving using wormhole-william
type Client struct {
	wormhole.Client

	// Save a reference to the window to avoid creating a new one when sending and receiving text
	display *textDisplay

	// Save a reference to the zip handler to avoid creating a new one each time when unzipping folders
	Zip *archiver.Zip

	// Notification holds the settings value for if we have notifications enabled or not.
	Notifications bool

	// DownloadPath holds the download path used for saving recvieved files.
	DownloadPath string
}

// ShowNotification sends a notification if c.Notifications is true.
func (c *Client) ShowNotification(title, content string) {
	if c.Notifications {
		fyne.CurrentApp().SendNotification(&fyne.Notification{Title: title, Content: content})
	}
}

// NewClient returns a new client for sending and receiving using wormhole-william
func NewClient() *Client {
	return &Client{display: createTextWindow(), Zip: &archiver.Zip{MkdirAll: true}}
}

// UserDownloadsFolder returns the downloads folder corresponding to the current user.
func UserDownloadsFolder() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		fyne.LogError("Could not get home dir", err)
	}

	return filepath.Join(dir, "Downloads")
}
