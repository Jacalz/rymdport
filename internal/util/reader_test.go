package util

import (
	"bytes"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestProgressReader(t *testing.T) {
	buffer := bytes.NewBufferString("content")
	size := int64(buffer.Len())

	var value, total int64
	teeReader := NewProgressReader(buffer, func(delta, max int64) {
		value += delta
		total = max
	}, size)

	temp := [1]byte{}
	for i := int64(0); i < size; i++ {
		_, err := teeReader.Read(temp[:])
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, int64(i+1), value)
		assert.Equal(t, size, total)
	}

	assert.Equal(t, size, value)
	assert.Equal(t, size, total)
}
