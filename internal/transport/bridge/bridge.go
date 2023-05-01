// Package bridge serves as a bridge between the transport backend and the Fyne ui
package bridge

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skip2/go-qrcode"
)

func newCodeDisplay(window fyne.Window) *fyne.Container {
	codeLabel := &widget.Label{Text: "Waiting for code..."}
	copyButton := &widget.Button{Icon: theme.ContentCopyIcon(), Importance: widget.LowImportance}
	clipboard := window.Clipboard()
	copyButton.OnTapped = func() {
		if codeLabel.Text != "Waiting for code..." {
			copyButton.SetIcon(theme.ConfirmIcon())
			clipboard.SetContent(codeLabel.Text)
		} else {
			copyButton.SetIcon(theme.CancelIcon())
		}

		time.Sleep(500 * time.Millisecond)
		copyButton.SetIcon(theme.ContentCopyIcon())
	}

	qrcodeButton := &widget.Button{
		Icon:       theme.InfoIcon(),
		Importance: widget.LowImportance,
		OnTapped: func() {
			if codeLabel.Text == "Waiting for code..." {
				return
			}

			code, err := qrcode.New("wormhole-transfer:"+codeLabel.Text, qrcode.High)
			if err != nil {
				fyne.LogError("Failed to encode qr code", err)
				return
			}

			code.BackgroundColor = theme.OverlayBackgroundColor()
			code.ForegroundColor = theme.ForegroundColor()

			qrcode := canvas.NewImageFromImage(code.Image(256))
			qrcode.SetMinSize(fyne.NewSize(200, 200))

			dialog.ShowCustom("QR Code", "Close", qrcode, window)
		},
	}

	return container.NewHBox(codeLabel, copyButton, qrcodeButton)
}
