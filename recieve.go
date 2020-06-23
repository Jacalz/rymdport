package main

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/psanford/wormhole-william/wormhole"
)

// Regular expression for verifying sync code.
var validCode = regexp.MustCompile(`^\d\d?(-\w{2,12}){2,6}$`)

func recieveData(code string, s settings, w fyne.Window) error {
	c := wormhole.Client{PassPhraseComponentLength: s.ComponentLength}

	msg, err := c.Receive(context.Background(), code)
	if err != nil {
		fyne.LogError("Error on receiving data", err)
		return err
	}

	text, err := ioutil.ReadAll(msg)
	if err != nil {
		fyne.LogError("Error on reading received data", err)
		return err
	}

	if msg.Type == wormhole.TransferText {
		textEntry := widget.NewMultiLineEntry()
		textEntry.SetText(string(text))

		dialog.ShowCustom("Received text", "Close", textEntry, w)
		return nil
	}

	f, err := os.Create(path.Join(s.DownloadPath, msg.Name))
	if err != nil {
		fyne.LogError("Error on creating file", err)
		return err
	}

	_, err = f.Write(text)
	if err != nil {
		if err2 := f.Close(); err2 != nil {
			fyne.LogError("Error on writing and closing the file", err)
			return err
		}

		fyne.LogError("Error on writing data to the file", err)
		return err

	}

	return f.Close()
}

func (s *settings) recieveTab(w fyne.Window) *widget.TabItem {
	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Enter code")

	codeButton := widget.NewButtonWithIcon("Download", theme.MoveDownIcon(), func() {
		code := codeEntry.Text
		if validCode.MatchString(code) {
			go func() {
				err := recieveData(code, *s, w)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				dialog.ShowInformation("Successful download", "The download completed without errors.", w)
			}()

			codeEntry.SetText("")
		}
	})

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), codeEntry, codeButton)

	recieveContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), codeContainer)

	return widget.NewTabItemWithIcon("Receive", theme.MoveDownIcon(), recieveContent)
}
