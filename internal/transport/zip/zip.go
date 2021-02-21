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

// Extract takes a reader and the length and then extracts it to the target.
func Extract(source io.ReaderAt, length int64, target string) error {
	reader, err := zip.NewReader(source, length)
	if err != nil {
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
		fyne.LogError("Dangerous filename detected", nil)
		return errors.New("dangerous filename detected: " + path)
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
