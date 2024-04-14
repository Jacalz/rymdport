package transport

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/Jacalz/rymdport/v3/zip"
	"github.com/rymdport/wormhole/wormhole"
)

// NewFileSend takes the chosen file and sends it using wormhole-william.
func (c *Client) NewFileSend(file fyne.URIReadCloser, progress wormhole.SendOption, code string) (string, chan wormhole.SendResult, error) {
	return c.SendFile(context.Background(), file.URI().Name(), file.(io.ReadSeeker), progress, wormhole.WithCode(code))
}

// NewDirSend takes a listable URI and sends it using wormhole-william.
func (c *Client) NewDirSend(dir fyne.ListableURI, progress wormhole.SendOption, code string) (string, chan wormhole.SendResult, error) {
	prefixStr, _ := filepath.Split(dir.Path())
	prefix := len(prefixStr) // Where the prefix ends. Doing it this way is faster and works when paths don't use same separator (\ or /).

	var files []wormhole.DirectoryEntry
	if err := filepath.Walk(dir.Path(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}

		files = append(files, wormhole.DirectoryEntry{
			Path: path[prefix:], // Instead of strings.TrimPrefix. Paths don't need match (e.g. "C:/home/dir" == "C:\home\dir").
			Mode: info.Mode(),
			Reader: func() (io.ReadCloser, error) {
				return os.Open(path) // #nosec - path is already cleaned by filepath.Walk
			},
		})

		return nil
	}); err != nil {
		return "", nil, err
	}

	return c.SendDirectory(context.Background(), dir.Name(), files, progress, wormhole.WithCode(code))
}

// NewMultipleFileSend sends multiple files as a directory send using wormhole-william.
func (c *Client) NewMultipleFileSend(uris []fyne.URI, progress wormhole.SendOption, code string) (string, chan wormhole.SendResult, error) {
	// We want to prefix all directory entires with the parent folder only.
	baseDir := filepath.Dir(uris[0].Path())
	subDir, dirName := filepath.Split(baseDir)
	prefix := len(subDir)

	files := make([]wormhole.DirectoryEntry, 0, len(uris))
	for _, uri := range uris {
		absPath, err := filepath.Abs(filepath.Join(baseDir, uri.Name()))
		if err != nil {
			return "", nil, err
		}

		if !strings.HasPrefix(absPath, baseDir) {
			return "", nil, zip.ErrorDangerousFilename
		}

		path := uri.Path()
		info, err := os.Lstat(path)
		if err != nil {
			return "", nil, err
		}

		if !info.IsDir() {
			files = append(files, wormhole.DirectoryEntry{
				Path: path[prefix:], // Instead of strings.TrimPrefix. Paths don't need match separators (e.g. "C:/home/dir" == "C:\home\dir").
				Mode: info.Mode(),
				Reader: func() (io.ReadCloser, error) {
					return os.Open(path) // #nosec - path is already cleaned by filepath.Walk
				},
			})

			continue
		}

		if err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			} else if info.IsDir() || !info.Mode().IsRegular() {
				return nil
			}

			files = append(files, wormhole.DirectoryEntry{
				Path: path[prefix:], // Instead of strings.TrimPrefix. Paths don't need match separators (e.g. "C:/home/dir" == "C:\home\dir").
				Mode: info.Mode(),
				Reader: func() (io.ReadCloser, error) {
					return os.Open(path) // #nosec - path is already cleaned by filepath.Walk
				},
			})

			return nil
		}); err != nil {
			return "", nil, err
		}
	}

	return c.SendDirectory(context.Background(), dirName, files, progress, wormhole.WithCode(code))
}

// NewTextSend takes a text input and sends the text using wormhole-william.
func (c *Client) NewTextSend(text string, progress wormhole.SendOption, code string) (string, chan wormhole.SendResult, error) {
	return c.SendText(context.Background(), text, progress, wormhole.WithCode(code))
}
