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

// SendFile takes the chosen file and sends it using wormhole-william.
func (b *Bridge) SendFile(file fyne.URIReadCloser, code chan string, progress wormhole.SendOption) error {
	c := wormhole.Client{PassPhraseComponentLength: b.ComponentLength}

	defer file.Close()

	f, err := os.Open(file.URI().String()[7:])
	if err != nil {
		fyne.LogError("Error on opening file", err)
		return err
	}

	defer f.Close() // #nosec - We are not writing to the file.

	codestr, status, err := c.SendFile(context.Background(), file.Name(), f, progress)
	if err != nil {
		fyne.LogError("Error on sending file", err)
		return err
	}

	code <- codestr

	if stat := <-status; stat.Error != nil {
		fyne.LogError("Error on status of share", err)
		return err
	} else if stat.OK {
		return nil
	}

	return nil
}

// SendDir takes a listable URI and sends it using wormhole-william.
func (b *Bridge) SendDir(dir fyne.ListableURI, code chan string, progress wormhole.SendOption) error {
	c := wormhole.Client{PassPhraseComponentLength: b.ComponentLength}

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
		return err
	}

	codestr, status, err := c.SendDirectory(context.Background(), dir.Name(), files, progress)
	if err != nil {
		fyne.LogError("Error on sending directory", err)
		return err
	}

	code <- codestr

	if stat := <-status; stat.Error != nil {
		fyne.LogError("Error on status of share", err)
		return err
	} else if stat.OK {
		return nil
	}

	return nil
}

// SendText takes a text input and sends the text using wormhole-william.
func (b *Bridge) SendText(text string, code chan string) error {
	c := wormhole.Client{PassPhraseComponentLength: b.ComponentLength}

	codestr, status, err := c.SendText(context.Background(), text) // TODO: Check why progress doesn't work for sending text.
	if err != nil {
		fyne.LogError("Error on sending text", err)
		return err
	}

	code <- codestr

	if stat := <-status; stat.Error != nil {
		fyne.LogError("Error on status of share", err)
		return err
	} else if stat.OK {
		return nil
	}

	return nil
}
