package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/psanford/wormhole-william/wormhole"
)

// Regular expression for verifying sync code.
var validCode = regexp.MustCompile(`^\d\d?-\w{2,12}-\w{2,12}$`)

func recieveData(code string, s settings) error {
	c := wormhole.Client{PassPhraseComponentLength: s.PassPhraseComponentLength}

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

	f, err := os.Create(path.Join(s.DownloadPath, msg.Name))
	if err != nil {
		fyne.LogError("Error on creating file", err)
		return err
	}

	_, err = f.Write(text)
	if err2 := f.Close(); err != nil || err2 != nil {
		errs := fmt.Errorf("write error: %v, file close error: %v", err, err2)
		fyne.LogError("Error on writing data to the file or closing the file", errs)
		return errs
	}

	return f.Close()
}

func (s *settings) recieveTab() *widget.TabItem {
	codeEntry := widget.NewEntry()
	codeEntry.SetPlaceHolder("Enter code")

	codeButton := widget.NewButtonWithIcon("Download", theme.MoveDownIcon(), func() {
		if validCode.MatchString(codeEntry.Text) {
			recieveData(codeEntry.Text, *s)
		}
	})

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2), codeEntry, codeButton)

	recieveContent := fyne.NewContainerWithLayout(layout.NewVBoxLayout(), codeContainer)

	return widget.NewTabItemWithIcon("Receive", theme.MoveDownIcon(), recieveContent)
}
