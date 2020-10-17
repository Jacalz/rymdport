package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/data/validation"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/Jacalz/wormhole-gui/internal/bridge/widgets"
)

func (ad *appData) recieveTab() *container.TabItem {
	codeEntry := widgets.NewPressEntry("Enter code")
	codeEntry.Validator = validation.NewRegexp(`^\d\d?(-\w{2,12}){2,6}$`, "Invalid code")

	recvList := widgets.NewRecvList()

	codeButton := widget.NewButtonWithIcon("Download", theme.DownloadIcon(), func() {
		go func(code string) {
			if err := codeEntry.Validate(); err == nil {
				filename, status := recvList.NewRecvItem()
				codeEntry.SetText("")

				go func() {
					err := ad.Bridge.RecieveData(code, filename, ad.App)
					if err != nil {
						status <- "Failed"
						dialog.ShowError(err, ad.Window)
						return
					}

					status <- "Completed"
					dialog.ShowInformation("Successful download", "The download completed without errors.", ad.Window)

					if ad.Notifications {
						ad.App.SendNotification(fyne.NewNotification("Receive completed", "The receive completed successfully"))
					}
				}()
			}
		}(codeEntry.Text)
	})
	codeEntry.OnReturn = codeButton.OnTapped

	box := widget.NewVBox(container.NewGridWithColumns(2, codeEntry, codeButton), widget.NewLabel(""))
	recvContent := fyne.NewContainerWithLayout(layout.NewBorderLayout(box, nil, nil, nil), box, recvList)

	return widget.NewTabItemWithIcon("Receive", theme.MoveDownIcon(), recvContent)
}
