// Package bridge serves as a bridge between the transport backend and the Fyne ui
package bridge

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// codeDisplay is a label extended to copy the code with a menu popup on rightclick.
type codeDisplay struct {
	widget.Label
	button    *widget.Button
	clipboard fyne.Clipboard
}

func (c *codeDisplay) copyOnPress() {
	if c.Text != "Waiting for code..." {
		c.button.SetIcon(theme.ConfirmIcon())
		c.clipboard.SetContent(c.Text)
	} else {
		c.button.SetIcon(theme.CancelIcon())
	}

	time.Sleep(500 * time.Millisecond)
	c.button.SetIcon(theme.ContentCopyIcon())
}

func newCodeDisplay() *fyne.Container {
	c := &codeDisplay{button: &widget.Button{Icon: theme.ContentCopyIcon(), Importance: widget.LowImportance},
		clipboard: fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()}
	c.ExtendBaseWidget(c)

	c.Text = "Waiting for code..."
	c.button.OnTapped = c.copyOnPress

	return container.NewHBox(c, c.button)
}
