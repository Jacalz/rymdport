// Package transport handles sending and receiving using wormhole-william
package transport

import (
	"fyne.io/fyne/v2"
	"github.com/rymdport/wormhole/wormhole"
)

// Client defines the client for handling sending and receiving using wormhole-william
type Client struct {
	wormhole.Client

	// App is a reference to the currently running Fyne application.
	App fyne.App

	// AskOnFileSave defines if we should ask where to save files or not.
	AskOnFileSave bool

	// CustomCode defines if we should pass a custom code or let wormhole-william generate on for us.
	CustomCode bool

	// DownloadPath holds the download path used for saving received files.
	DownloadPath string

	// NoExtractDirectory specifies if we should extract directories or not.
	NoExtractDirectory bool

	// Notification holds the settings value for if we have notifications enabled or not.
	Notifications bool

	// OverwriteExisting holds the settings value for if we should overwrite already existing files.
	OverwriteExisting bool
}

// ShowNotification sends a notification if c.Notifications is true.
func (c *Client) ShowNotification(title, content string) {
	if c.Notifications {
		c.App.SendNotification(&fyne.Notification{Title: title, Content: content})
	}
}

// NewClient returns a new client for sending and receiving using wormhole-william
func NewClient(app fyne.App) *Client {
	return &Client{App: app}
}
