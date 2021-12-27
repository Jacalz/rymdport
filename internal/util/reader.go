package util

import (
	"io"
)

type teeReaderAt struct {
	r io.ReaderAt
	w io.Writer
}

// ReadAt wraps the reader and calls write on it.
func (t *teeReaderAt) ReadAt(p []byte, off int64) (int, error) {
	n, err := t.r.ReadAt(p, off)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}

	return n, err
}

// TeeReaderAt returns a wrapped ReaderAt that writes what is being read.
func TeeReaderAt(r io.ReaderAt, w io.Writer) io.ReaderAt {
	return &teeReaderAt{r, w}
}
