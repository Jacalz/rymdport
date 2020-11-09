package ui

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/data/validation"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/internal/bridge"
	"github.com/Jacalz/wormhole-gui/internal/bridge/widgets"
)

type recv struct {
	codeEntry  *widgets.PressEntry
	codeButton *widget.Button

	recvList *widgets.RecvList

	bridge      *bridge.Bridge
	appSettings *AppSettings
	window      fyne.Window
	app         fyne.App
}

func newRecv(a fyne.App, w fyne.Window, b *bridge.Bridge, as *AppSettings) *recv {
	return &recv{app: a, window: w, bridge: b, appSettings: as}
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
	r.codeEntry = widgets.NewPressEntry("Enter code")
	r.codeEntry.Validator = validation.NewRegexp(`^\d\d?(-\w{2,12}){2,6}$`, "Invalid code")
	r.codeEntry.OnReturn = r.onRecv

	r.codeButton = &widget.Button{Text: "Download", Icon: theme.DownloadIcon(), OnTapped: r.onRecv}

	r.recvList = widgets.NewRecvList(r.bridge)

	box := container.NewVBox(container.NewGridWithColumns(2, r.codeEntry, r.codeButton), &widget.Label{})
	return container.NewBorder(box, nil, nil, nil, r.recvList)
}

func (r *recv) tabItem() *container.TabItem {
	return container.NewTabItemWithIcon("Receive", theme.MailSendIcon(), r.buildUI())
}
