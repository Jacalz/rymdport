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
	"github.com/rymdport/wormhole/wormhole"
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
func (c *Client) NewReceive(code string) (*wormhole.IncomingMessage, error) {
	msg, err := c.Receive(context.Background(), code)
	if err != nil {
		fyne.LogError("Error on receiving data", err)
		return nil, bail(msg, err)
	}

	return msg, nil
}

// SaveToDisk saves the incomming file or directory transfer to the disk.
func (c *Client) SaveToDisk(msg *wormhole.IncomingMessage, targetPath string, progress func(int64, int64)) (err error) {
	contents := util.NewProgressReader(msg, progress, msg.TransferBytes)

	if msg.Type == wormhole.TransferDirectory && c.NoExtractDirectory {
		targetPath += ".zip" // We are saving the zip-file and not extracting.
	}

	if !c.OverwriteExisting {
		_, err := os.Stat(targetPath)
		if err == nil || os.IsExist(err) {
			targetPath, err = addFileIncrement(targetPath)
			if err != nil {
				fyne.LogError("Error on trying to create non-duplicate filename", err)
				return bail(msg, err)
			}
		}
	}

	if msg.Type == wormhole.TransferFile || c.NoExtractDirectory {
		return writeToFile(targetPath, msg, contents)
	}

	// We are reading the transferred bytes twice. First from msg to temp file and then from temp.
	contents.Max *= 2

	return writeToDirectory(targetPath, msg, contents, progress)
}

func writeToDirectory(targetPath string, msg *wormhole.IncomingMessage, contents util.ProgressReader, progress func(int64, int64)) (err error) {
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

	var n int64
	n, err = io.Copy(tmp, contents)
	if err != nil {
		fyne.LogError("Error on copying contents to file", err)
		return
	}

	err = zip.ExtractSafe(util.NewProgressReaderAt(tmp, progress, contents.Max),
		n, targetPath, msg.UncompressedBytes, msg.FileCount)
	if err != nil {
		fyne.LogError("Error on unzipping contents", err)
		return
	}

	progress(0, 1) // Workaround for progress sometimes stopping at 99%.

	return
}

func writeToFile(destination string, msg *wormhole.IncomingMessage, contents util.ProgressReader) (err error) {
	file, err := os.Create(destination) // #nosec Path is cleaned by filepath.Join().
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
		return
	}

	return
}

// addFileIncrement tries to add a number to the end of the filename if a duplicate exists.
// If it fails to do so after five tries, it till return the given path and an error.
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

	return path, errorTooManyDuplicates
}
