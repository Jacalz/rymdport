package util

import "io"

// ProgressReader is a reader wrapper that calls a function for each part content being read.
type ProgressReader struct {
	readat io.ReaderAt
	read   io.Reader

	Max      int64
	progress func(delta int64, max int64)
}

// ReadAt wraps the ReaderAt and updates the progress bar.
func (p *ProgressReader) ReadAt(buf []byte, off int64) (int, error) {
	n, err := p.readat.ReadAt(buf, off)
	p.progress(int64(n), p.Max)
	return n, err
}

// Read wraps the Reader and updates the progress bar.
func (p *ProgressReader) Read(buf []byte) (int, error) {
	n, err := p.read.Read(buf)
	p.progress(int64(n), p.Max)
	return n, err
}

// NewProgressReaderAt returns a wrapped ReaderAt that updates the progress bar.
func NewProgressReaderAt(r io.ReaderAt, p func(delta int64, max int64), max int64) *ProgressReader {
	return &ProgressReader{readat: r, progress: p, Max: max}
}

// NewProgressReader returns a wrapped Reader that updates the progress bar.
func NewProgressReader(r io.Reader, p func(delta int64, max int64), max int64) *ProgressReader {
	return &ProgressReader{read: r, progress: p, Max: max}
}
