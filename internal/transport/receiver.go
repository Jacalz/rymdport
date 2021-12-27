package transport

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"github.com/Jacalz/wormhole-gui/v2/internal/transport/zip"
	"github.com/Jacalz/wormhole-gui/v2/internal/util"
	"github.com/psanford/wormhole-william/wormhole"
)

func bail(msg *wormhole.IncomingMessage, err error) error {
	if msg == nil || msg.Type == wormhole.TransferText { // Rejecting text receives is not possible.
		return err
	} else if rerr := msg.Reject(); rerr != nil {
		return rerr
	}

	return err
}

// NewReceive runs a receive using wormhole-william and handles types accordingly.
func (c *Client) NewReceive(code string, pathname chan string, progress *util.ProgressBar) (err error) {
	msg, err := c.Receive(context.Background(), code)
	if err != nil {
		pathname <- "fail" // We want to always send a URI, even on fail, in order to not block goroutines.
		fyne.LogError("Error on receiving data", err)
		return bail(msg, err)
	}

	progress.Max = float64(msg.TransferBytes64)
	contents := io.TeeReader(msg, progress)

	if msg.Type == wormhole.TransferText {
		pathname <- "text"

		text := &bytes.Buffer{}
		text.Grow(int(msg.TransferBytes64))

		_, err := io.Copy(text, contents)
		if err != nil {
			fyne.LogError("Could not copy the received text", err)
			return err
		}

		c.showTextReceiveWindow(text)
		return nil
	}

	path := filepath.Join(c.DownloadPath, msg.Name)
	pathname <- path

	if !c.OverwriteExisting {
		if _, err := os.Stat(path); err == nil || os.IsExist(err) {
			fyne.LogError("Settings prevent overwriting existing files and folders", err)
			return bail(msg, os.ErrExist)
		}
	}

	if msg.Type == wormhole.TransferFile {
		var file *os.File
		file, err = os.Create(path)
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

		_, err = io.Copy(file, contents)
		if err != nil {
			fyne.LogError("Error on copying contents to file", err)
			return err
		}

		return
	}

	// We are reading the data twice. First from msg to temp file and then from temp.
	progress.Max += float64(msg.TransferBytes64)

	tmp, err := ioutil.TempFile("", msg.Name+"-*.zip.tmp")
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

	n, err := io.Copy(tmp, contents)
	if err != nil {
		fyne.LogError("Error on copying contents to file", err)
		return err
	}

	err = zip.Extract(util.TeeReaderAt(tmp, progress), n, path)
	if err != nil {
		fyne.LogError("Error on unzipping contents", err)
		return err
	}

	// TODO: Progress sometimes stops at 99%. Is it the offset that isn't accounted for?
	progress.Done()

	return
}
