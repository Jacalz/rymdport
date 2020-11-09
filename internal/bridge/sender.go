package bridge

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne"
	"github.com/psanford/wormhole-william/wormhole"
)

// NewFileSend takes the chosen file and sends it using wormhole-william.
func (b *Bridge) NewFileSend(file fyne.URIReadCloser, progress wormhole.SendOption) (string, chan wormhole.SendResult, *os.File, error) {
	err := file.Close() // Not currently used due to not being an io.ReadSeeker
	if err != nil {
		fyne.LogError("Error on closing file", err)
	}

	f, err := os.Open(file.URI().String()[7:])
	if err != nil {
		fyne.LogError("Error on opening file", err)
		return "", nil, nil, f.Close()
	}

	code, result, err := b.SendFile(context.Background(), file.URI().Name(), f, progress)
	return code, result, f, err
}

// NewDirSend takes a listable URI and sends it using wormhole-william.
func (b *Bridge) NewDirSend(dir fyne.ListableURI, progress wormhole.SendOption) (string, chan wormhole.SendResult, error) {
	dirpath := dir.String()[7:]
	prefix, _ := filepath.Split(dirpath)

	var files []wormhole.DirectoryEntry
	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		files = append(files, wormhole.DirectoryEntry{
			Path: strings.TrimPrefix(path, prefix),
			Mode: info.Mode(),
			Reader: func() (io.ReadCloser, error) {
				return os.Open(filepath.Clean(path))
			},
		})

		return nil
	})

	if err != nil {
		fyne.LogError("Error on walking directory", err)
		return "", nil, err
	}

	return b.SendDirectory(context.Background(), dir.Name(), files, progress)
}

// NewTextSend takes a text input and sends the text using wormhole-william.
func (b *Bridge) NewTextSend(text string, progress wormhole.SendOption) (string, chan wormhole.SendResult, error) {
	return b.SendText(context.Background(), text, progress)
}
