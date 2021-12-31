package util

import (
	"io"
)

type teeReader struct {
	readat io.ReaderAt
	read   io.Reader
	prog   *ProgressBar
}

// ReadAt wraps the ReaderAt and updates the progress bar.
func (t *teeReader) ReadAt(p []byte, off int64) (int, error) {
	n, err := t.readat.ReadAt(p, off)
	t.prog.SetValue(t.prog.Value + float64(n))
	return n, err
}

// Read wraps the Reader and updates the progress bar.
func (t *teeReader) Read(p []byte) (int, error) {
	n, err := t.read.Read(p)
	t.prog.SetValue(t.prog.Value + float64(n))
	return n, err
}

// TeeReaderAt returns a wrapped ReaderAt that updates the progress bar.
func TeeReaderAt(r io.ReaderAt, p *ProgressBar) io.ReaderAt {
	return &teeReader{readat: r, prog: p}
}

// TeeReader returns a wrapped Reader that updates the progress bar.
func TeeReader(r io.Reader, p *ProgressBar) io.Reader {
	return &teeReader{read: r, prog: p}
}
