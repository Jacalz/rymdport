package zip

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/klauspost/compress/zip"
)

var errorDangerousFilename = errors.New("dangerous filename detected")

// Extract takes a reader and the length and then extracts it to the target.
func Extract(source io.ReaderAt, length int64, target string) error {
	reader, err := zip.NewReader(source, length)
	if err != nil {
		fyne.LogError("Could not create zip reader", err)
		return err
	}

	for _, file := range reader.File {
		if err := extractFile(file, target); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file *zip.File, target string) (err error) {
	path, err := filepath.Abs(filepath.Join(target, file.Name))
	if err != nil {
		fyne.LogError("Could not calculate the ABS path", err)
		return err
	}

	if !strings.HasPrefix(path, target) {
		fyne.LogError("Dangerous filename detected", errorDangerousFilename)
		return errorDangerousFilename
	}

	fileReader, err := file.Open()
	if err != nil {
		fyne.LogError("Could not open the zip file", err)
		return err
	}

	defer func() {
		if cerr := fileReader.Close(); cerr != nil {
			fyne.LogError("Could not close the zip file reader", err)
			err = cerr
		}
	}()

	if file.FileInfo().IsDir() {
		err = os.MkdirAll(path, 0750)
		if err != nil {
			fyne.LogError("Could not create the directory", err)
			return err
		}

		return nil
	}

	targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode()) // #nosec - The path has already been cleaned by filepath.Abs()
	if err != nil {
		fyne.LogError("Could not create the target file", err)
		return err
	}

	defer func() {
		if cerr := targetFile.Close(); cerr != nil {
			fyne.LogError("Could not close the target file", err)
			err = cerr
		}
	}()

	_, err = io.Copy(targetFile, fileReader)
	if err != nil {
		fyne.LogError("Could not copy the contents", err)
		return err
	}

	return
}
