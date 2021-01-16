package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/data/validation"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/transport"
	"github.com/Jacalz/wormhole-gui/internal/transport/bridge"
)

type recv struct {
	codeEntry  *widget.Entry
	codeButton *widget.Button

	recvList *bridge.RecvList

	bridge      *transport.Client
	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newRecv(a fyne.App, w fyne.Window, c *transport.Client, as *AppSettings) *recv {
	return &recv{app: a, window: w, bridge: c, appSettings: as}
}

func (r *recv) onRecv() {
	if err := r.codeEntry.Validate(); err != nil {
		dialog.ShowInformation("Invalid code", "The code is invalid. Please try again.", r.window)
		return
	}

	r.recvList.NewReceive(r.codeEntry.Text)
	r.codeEntry.SetText("")
}

func (r *recv) buildUI() *fyne.Container {
	r.codeEntry = &widget.Entry{PlaceHolder: "Enter code", OnSubmitted: func(_ string) { r.onRecv() },
		Validator: validation.NewRegexp(`^\d\d?(-\w{2,12}){2,6}$`, "The code is invalid"),
	}

	r.codeButton = &widget.Button{Text: "Download", Icon: theme.DownloadIcon(), OnTapped: r.onRecv}

	r.recvList = bridge.NewRecvList(r.bridge)

	box := container.NewVBox(container.NewGridWithColumns(2, r.codeEntry, r.codeButton), &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, r.recvList)
}

func (r *recv) tabItem() *container.TabItem {
	return &container.TabItem{Text: "Receive", Icon: theme.DownloadIcon(), Content: r.buildUI()}
}
