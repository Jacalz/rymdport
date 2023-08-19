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
	codeLabel := &widget.Label{Text: "Waiting for code...", Truncation: fyne.TextTruncateEllipsis}
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

	return container.New(codeLayout{}, codeLabel, copyButton)
}

type codeLayout struct{}

func (c codeLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	displacement := size.Width * 0.8

	objects[0].Move(fyne.NewSquareOffsetPos(0))
	objects[0].Resize(fyne.NewSize(displacement, size.Height))

	objects[1].Move(fyne.NewPos(displacement, 0))
	objects[1].Resize(fyne.NewSquareSize(size.Height))

}

func (c codeLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	leftMin := objects[0].MinSize()
	rightMin := objects[1].MinSize()

	return fyne.NewSize(leftMin.Width+leftMin.Width, fyne.Max(leftMin.Height, rightMin.Height))
}
