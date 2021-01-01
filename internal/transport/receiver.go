package transport

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
	"github.com/psanford/wormhole-william/wormhole"
)

func bail(msg *wormhole.IncomingMessage, err error) error {
	if rerr := msg.Reject(); rerr != nil {
		return rerr
	}

	return err
}

// NewReceive runs a receive using wormhole-william and handles types accordingly.
func (c *Client) NewReceive(code string, uri chan fyne.URI) error {
	msg, err := c.Receive(context.Background(), code)
	if err != nil {
		fyne.LogError("Error on receiving data", err)
		return bail(msg, err)
	}

	if msg.Type == wormhole.TransferText {
		content, err := ioutil.ReadAll(msg)
		if err != nil {
			fyne.LogError("Error on reading received data", err)
			return err
		}

		c.displayReceivedText(content)

		uri <- storage.NewURI("Text Snippet")
	}

	path := filepath.Join(c.DownloadPath, msg.Name)

	switch msg.Type {
	case wormhole.TransferFile:
		if !c.Zip.OverwriteExisting {
			if _, err := os.Stat(path); err == nil || os.IsExist(err) {
				fyne.LogError("Error on creating file, settings prevent overwriting existing files", err)
				return bail(msg, os.ErrExist)
			}
		}

		file, err := os.Create(path)
		if err != nil {
			fyne.LogError("Error on creating file", err)
			return bail(msg, err)
		}

		defer func() {
			if cerr := file.Close(); cerr != nil {
				fyne.LogError("Error on closing file", err)
				err = cerr
			}
		}()

		_, err = io.Copy(file, ioutil.NopCloser(msg))
		if err != nil {
			fyne.LogError("Error on copying contents to file", err)
			return err
		}

		uri <- storage.NewFileURI(path)
	case wormhole.TransferDirectory:
		tmp, err := ioutil.TempFile("", msg.Name+".zip.tmp")
		if err != nil {
			fyne.LogError("Error on creating tempfile", err)
			return bail(msg, err)
		}

		defer func() {
			if cerr := tmp.Close(); cerr != nil {
				fyne.LogError("Error on closing file", err)
				err = cerr
			}

			if rerr := os.Remove(tmp.Name()); rerr != nil {
				fyne.LogError("Error on removing temp file", err)
				err = rerr
			}
		}()

		_, err = io.Copy(tmp, ioutil.NopCloser(msg))
		if err != nil {
			fyne.LogError("Error on copying contents to file", err)
			return err
		}

		err = c.Zip.Unarchive(tmp.Name(), path)
		if err != nil {
			fyne.LogError("Error on unzipping contents", err)
			return err
		}

		uri <- storage.NewFileURI(path)
	}

	return nil
}
