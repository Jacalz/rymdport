// Package transport handles sending and receiving using wormhole-william
package transport

import (
	"fyne.io/fyne/v2"
	"github.com/psanford/wormhole-william/wormhole"
)

// Client defines the client for handling sending and receiving using wormhole-william
type Client struct {
	wormhole.Client

	app fyne.App

	// mainWindow holds a reference to the main window so we can focus it after a send.
	mainWindow fyne.Window

	// Save a reference to the window to avoid creating a new one when sending and receiving text.
	textSendWindow *textSendWindow
	textRecvWindow *textRecvWindow

	// Notification holds the settings value for if we have notifications enabled or not.
	Notifications bool

	// OverwriteExisting holds the settings value for if we should overwrite already existing files.
	OverwriteExisting bool

	// DownloadPath holds the download path used for saving received files.
	DownloadPath string

	// Defines if we should pass a custom code or let wormhole-william generate on for us.
	CustomCode bool
}

// ShowNotification sends a notification if c.Notifications is true.
func (c *Client) ShowNotification(title, content string) {
	if c.Notifications {
		c.app.SendNotification(&fyne.Notification{Title: title, Content: content})
	}
}

// NewClient returns a new client for sending and receiving using wormhole-william
func NewClient(app fyne.App, mainWindow fyne.Window) *Client {
	return &Client{app: app, mainWindow: mainWindow}
}
