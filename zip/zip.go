// Package zip contains an implementation of a zip extractor.
package zip

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
)

var (
	// ErrorDangerousFilename indicates that a dangerous filename was found.
	ErrorDangerousFilename = errors.New("dangerous filename detected")

	// ErrorSizeMismatch indicates that the uncompressed size was unexpected.
	ErrorSizeMismatch = errors.New("mismatch between offered and actual size")

	// ErrorFileCountMismatch indicates that the file count was unexpected.
	ErrorFileCountMismatch = errors.New("mismatch between offered and actual file count")
)

// ExtractSafe works like Extract() but verifies that the uncompressed size and file count are as expected.
// This can only be used if you know the file count and uncompressed size before extracting.
func ExtractSafe(source io.ReaderAt, length int64, target string, uncompressedBytes int64, files int) error {
	reader, err := zip.NewReader(source, length)
	if err != nil {
		return err
	}

	if files < len(reader.File) {
		return ErrorFileCountMismatch
	}

	actualUncompressedSize := uint64(0)
	for _, f := range reader.File {
		actualUncompressedSize += f.FileHeader.UncompressedSize64
	}
	if uncompressedBytes < 0 || actualUncompressedSize > uint64(uncompressedBytes) {
		return ErrorSizeMismatch
	}

	return extractFiles(reader, target)
}

// Extract takes a reader and the length and then extracts it to the target.
// The target should be the path to a folder where the extracted files can be placed.
func Extract(source io.ReaderAt, length int64, target string) error {
	reader, err := zip.NewReader(source, length)
	if err != nil {
		return err
	}

	return extractFiles(reader, target)
}

func extractFiles(reader *zip.Reader, target string) error {
	for _, file := range reader.File {
		if err := extractFile(file, target); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(file *zip.File, target string) error {
	path, err := filepath.Abs(filepath.Join(target, file.Name))
	if err != nil {
		return err
	}

	if !strings.HasPrefix(path, target) {
		return ErrorDangerousFilename
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}

	defer fileReader.Close()

	if file.FileInfo().IsDir() {
		return os.MkdirAll(path, 0o750)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode()) // #nosec - The path has already been cleaned by filepath.Abs()
	if err != nil {
		return err
	}

	defer targetFile.Close()

	_, err = io.Copy(targetFile, fileReader)
	return err
}
