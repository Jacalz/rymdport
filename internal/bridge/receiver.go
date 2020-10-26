package bridge

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"fyne.io/fyne"
	"github.com/mholt/archiver/v3"
	"github.com/psanford/wormhole-william/wormhole"
)

// writeFile writes the given file to the file system and closes it when done.
func writeFile(file io.WriteCloser, content []byte) error {
	_, err := file.Write(content)
	if err != nil {
		if err2 := file.Close(); err2 != nil {
			fyne.LogError("Error on writing and closing the file", err)
			return err
		}

		fyne.LogError("Error on writing data to the file", err)
		return err
	}

	return file.Close() // Not defering close du to security issues when writing to a file.
}

// RecieveData runs a receive using wormhole-william and handles types accordingly.
func (b *Bridge) RecieveData(code string, fileName chan string, a fyne.App) error {
	c := wormhole.Client{PassPhraseComponentLength: b.ComponentLength}

	msg, err := c.Receive(context.Background(), code)
	if err != nil {
		fyne.LogError("Error on receiving data", err)
		return err
	}

	switch msg.Type {
	case wormhole.TransferText:
		content, err := ioutil.ReadAll(msg)
		if err != nil {
			fyne.LogError("Error on reading received data", err)
			return err
		}

		displayRecievedText(a, string(content))

		fileName <- "Text Snippet"
	case wormhole.TransferFile:
		file, err := os.Create(path.Join(b.DownloadPath, msg.Name))
		if err != nil {
			fyne.LogError("Error on creating file", err)
			return err
		}

		_, err = io.Copy(file, ioutil.NopCloser(msg))
		if err != nil {
			fyne.LogError("Error on copying contents to file", err)
			return err
		}

		fileName <- msg.Name
	case wormhole.TransferDirectory:
		dir := filepath.Join(b.DownloadPath, msg.Name)

		tmp, err := ioutil.TempFile("", msg.Name+".zip.tmp")
		if err != nil {
			fyne.LogError("Error on creating tempfile", err)
			return err
		}

		defer tmp.Close() // #nosec - We are not writing to the file
		defer os.Remove(tmp.Name())

		_, err = io.Copy(tmp, ioutil.NopCloser(msg))
		if err != nil {
			fyne.LogError("Error on copying contents to file", err)
			return err
		}

		err = archiver.NewZip().Unarchive(tmp.Name(), dir)
		if err != nil {
			fyne.LogError("Error on unzipping contents", err)
			return err
		}

		fileName <- msg.Name
	}

	return nil
}
