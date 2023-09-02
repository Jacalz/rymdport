package transport

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"fyne.io/fyne/v2"
	"github.com/Jacalz/rymdport/v3/internal/util"
	"github.com/Jacalz/rymdport/v3/zip"
	"github.com/psanford/wormhole-william/wormhole"
)

var errorTooManyDuplicates = errors.New("too many duplicates found. Stopped trying to add new numbers to end")

func bail(msg *wormhole.IncomingMessage, err error) error {
	if msg == nil || msg.Type == wormhole.TransferText { // Rejecting text receives is not possible.
		return err
	} else if rerr := msg.Reject(); rerr != nil {
		return rerr
	}

	return err
}

// NewReceive runs a receive using wormhole-william and handles types accordingly.
func (c *Client) NewReceive(code string, pathname chan string, progress func(int64, int64)) (err error) {
	msg, err := c.Receive(context.Background(), code)
	if err != nil {
		pathname <- "" // We want to always send a URI, even on fail, in order to not block goroutines.
		fyne.LogError("Error on receiving data", err)
		return bail(msg, err)
	}

	contents := util.NewProgressReader(msg, progress, msg.TransferBytes64)

	if msg.Type == wormhole.TransferText {
		pathname <- "Text Snippet"

		text := make([]byte, int(msg.TransferBytes64))
		_, err := io.ReadFull(contents, text)
		if err != nil {
			fyne.LogError("Could read the received text", err)
			return err
		}

		c.showTextReceiveWindow(text)
		return nil
	}

	path := filepath.Join(c.DownloadPath, msg.Name)
	pathname <- path

	if !c.OverwriteExisting {
		if _, err := os.Stat(path); err == nil || os.IsExist(err) {
			new, err := addFileIncrement(path)
			if err != nil {
				fyne.LogError("Error on trying to create non-duplicate filename", err)
				return bail(msg, err)
			}

			path = new
		}
	}

	if msg.Type == wormhole.TransferFile {
		var file *os.File
		file, err = os.Create(path) // #nosec Path is cleaned by filepath.Join().
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

	// We are reading the transferred bytes twice. First from msg to temp file and then from temp.
	contents.Max *= 2

	tmp, err := os.CreateTemp("", msg.Name+"-*.zip.tmp")
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

	err = zip.ExtractSafe(
		util.NewProgressReaderAt(tmp, progress, contents.Max),
		n, path, msg.UncompressedBytes64, msg.FileCount,
	)
	if err != nil {
		fyne.LogError("Error on unzipping contents", err)
		return err
	}

	progress(0, 1) // Workaround for progress sometimes stopping at 99%.

	return
}

func addFileIncrement(path string) (string, error) {
	base := filepath.Dir(path)
	ext := filepath.Ext(path)
	name := filepath.Base(path)
	name = name[:len(name)-len(ext)]

	// Add number at the end. Cap it to avoid doing checks forever.
	for i := 1; i <= 5; i++ {
		nr := strconv.Itoa(i)
		incremented := filepath.Join(base, name+"("+nr+")"+ext)

		_, err := os.Stat(incremented)
		if err != nil {
			if os.IsNotExist(err) {
				return incremented, nil
			}

			return "", err
		}
	}

	return "", errorTooManyDuplicates
}
