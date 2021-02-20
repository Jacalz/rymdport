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

// Extract extracts the zip archive at source to the target.
// If overwrite is set, it will overwrite existing files.
func Extract(source, target string) (err error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		fyne.LogError("Could not open the zip archive file", err)
		return err
	}

	defer func() {
		if cerr := reader.Close(); cerr != nil {
			fyne.LogError("Could not close the zip reader", cerr)
			err = cerr
		}
	}()

	for _, file := range reader.File {
		path, err := filepath.Abs(filepath.Join(target, file.Name))
		if err != nil {
			fyne.LogError("Could not calculate the ABS path", err)
			return err
		}

		if !strings.HasPrefix(path, target) {
			fyne.LogError("Dangerous filename detected", nil)
			return errors.New("dangerous filename detected: " + path)
		}

		fileReader, err := file.Open()
		if err != nil {
			fyne.LogError("Could not open the zip file", err)
			return err
		}

		err = os.MkdirAll(filepath.Dir(path), 0750)
		if err != nil {
			fyne.LogError("Could not create the directory", err)
			return err
		}

		targetFile, err := os.Create(path)
		if err != nil {
			fyne.LogError("Could not create the target file", err)
			return err
		}

		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			fyne.LogError("Could not copy the contents", err)
			return err
		}

		err = fileReader.Close()
		if err != nil {
			fyne.LogError("Could not close the zip file reader", err)
			return err
		}

		err = targetFile.Close()
		if err != nil {
			fyne.LogError("Could not close the target file", err)
			return err
		}
	}

	return
}
