package util

import (
	"io"
)

type TeeReader struct {
	readat io.ReaderAt
	read   io.Reader

	Max      int64
	progress func(delta int64, max int64)
}

// ReadAt wraps the ReaderAt and updates the progress bar.
func (t *TeeReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := t.readat.ReadAt(p, off)
	t.progress(int64(n), t.Max)
	return n, err
}

// Read wraps the Reader and updates the progress bar.
func (t *TeeReader) Read(p []byte) (int, error) {
	n, err := t.read.Read(p)
	t.progress(int64(n), t.Max)
	return n, err
}

// NewTeeReaderAt returns a wrapped ReaderAt that updates the progress bar.
func TeeReaderAt(r io.ReaderAt, p func(delta int64, max int64), max int64) *TeeReader {
	return &TeeReader{readat: r, progress: p, Max: max}
}

// NewTeeReader returns a wrapped Reader that updates the progress bar.
func NewTeeReader(r io.Reader, p func(delta int64, max int64), max int64) *TeeReader {
	return &TeeReader{read: r, progress: p, Max: max}
}
