package transport

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

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

func writeToDirectory(targetPath string, msg *wormhole.IncomingMessage, contents util.ProgressReader, progress func(int64, int64)) error {
	tmp, err := os.CreateTemp("", msg.Name+"-*.zip.tmp")
	if err != nil {
		return bail(msg, err)
	}

	defer tmp.Close()
	defer os.Remove(tmp.Name())

	var n int64
	n, err = io.Copy(tmp, contents)
	if err != nil {
		return err
	}

	err = zip.ExtractSafe(util.NewProgressReaderAt(tmp, progress, contents.Max),
		n, targetPath, msg.UncompressedBytes, msg.FileCount)
	if err != nil {
		return err
	}

	progress(0, 1) // Workaround for progress sometimes stopping at 99%.
	return nil
}

func writeToFile(destination string, msg *wormhole.IncomingMessage, contents util.ProgressReader) error {
	file, err := os.Create(destination) // #nosec Path is cleaned by filepath.Join().
	if err != nil {
		return bail(msg, err)
	}

	defer file.Close()

	_, err = io.Copy(file, contents)
	return err
}

// addFileIncrement tries to add a number to the end of the filename if a duplicate exists.
// If it fails to do so after five tries, it till return the given path and an error.
func addFileIncrement(path string) (string, error) {
	base := filepath.Dir(path)
	ext := filepath.Ext(path)
	name := filepath.Base(path)
	name = name[:len(name)-len(ext)]

	// Add number at the end. Cap it to avoid doing checks forever.
	for i := range 5 {
		nr := strconv.Itoa(i + 1)
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
