package main

import (
	"regexp"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/Jacalz/wormhole-gui/bridge/widgets"
)

// Regular expression for verifying sync code.
var validCode = regexp.MustCompile(`^\d\d?(-\w{2,12}){2,6}$`)

func (ad *appData) recieveTab() *widget.TabItem {
	codeEntry := widgets.NewPressEntry("Enter code")
	codeButton := widget.NewButtonWithIcon("Download", theme.DownloadIcon(), nil)

	recieveGrid := fyne.NewContainerWithLayout(layout.NewGridLayout(2), widgets.NewBoldLabel("Filename"), widgets.NewBoldLabel("Status"))

	codeButton.OnTapped = func() {
		go func() {
			code := codeEntry.Text
			if validCode.MatchString(code) {
				file := make(chan string)
				codeEntry.SetText("")

				filename := widget.NewLabel("Waiting for filename")
				recieveGrid.AddObject(filename)

				finished := widget.NewLabel("Waiting for status")
				recieveGrid.AddObject(finished)

				go func() {
					err := ad.Bridge.RecieveData(code, file, ad.App)
					if err != nil {
						finished.SetText("Failed")
						dialog.ShowError(err, ad.Window)
						return
					}

					finished.SetText("Completed")
					dialog.ShowInformation("Successful download", "The download completed without errors.", ad.Window)
					if ad.Notifications {
						ad.App.SendNotification(fyne.NewNotification("Receive completed", "The receive completed successfully"))
					}
				}()

				go filename.SetText(<-file)
			}
		}()
	}

	codeEntry.OnReturn = codeButton.OnTapped

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), codeEntry, codeButton)

	recieveContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), codeContainer, widget.NewLabel(""), recieveGrid)

	return widget.NewTabItemWithIcon("Receive", theme.MoveDownIcon(), recieveContent)
}
