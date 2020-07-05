package bridge

import (
	"context"
	"io/ioutil"
	"os"
	"path"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/psanford/wormhole-william/wormhole"
)

// RecieveData runs a receive using wormhole-william and handles types accordingly.
func (b *Bridge) RecieveData(code string, fileName chan string, w *fyne.Window) error {
	c := wormhole.Client{PassPhraseComponentLength: b.ComponentLength}

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
		fileName <- "Text Snippet"

		dialog.ShowCustom("Received text", "Close", textEntry, *w)
		return nil
	}

	fileName <- msg.Name

	f, err := os.Create(path.Join(b.DownloadPath, msg.Name))
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
