// Package bridge serves as a bridge between the transport backend and the Fyne ui
package bridge

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func newCodeDisplay(window fyne.Window) *fyne.Container {
	codeLabel := &widget.Label{Text: "Waiting for code..."}
	copyButton := &widget.Button{Icon: theme.ContentCopyIcon(), Importance: widget.LowImportance}

	copyButton.OnTapped = func() {
		if codeLabel.Text != "Waiting for code..." {
			copyButton.SetIcon(theme.ConfirmIcon())
			window.Clipboard().SetContent(codeLabel.Text)
		} else {
			copyButton.SetIcon(theme.CancelIcon())
		}

		time.Sleep(500 * time.Millisecond)
		copyButton.SetIcon(theme.ContentCopyIcon())
	}

	return container.NewHBox(codeLabel, copyButton)
}
