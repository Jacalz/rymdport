package util

import (
	"testing"

	"fyne.io/fyne/v2"
	"github.com/stretchr/testify/assert"
)

var globalValidationError error

func BenchmarkCodeValidator(b *testing.B) {
	local := error(nil)

	for i := 0; i < b.N; i++ {
		local = CodeValidator("125-upset-universe-mistake")
	}

	globalValidationError = local
}

func TestCodeValidator(t *testing.T) {
	validate := CodeValidator

	// This is so that the Entry doesn't report invalid when empty.
	valid := validate("")
	assert.NoError(t, valid)

	valid = validate("invalid-code")
	assert.Error(t, valid)

	valid = validate("IOI-")
	assert.Error(t, valid)

	valid = validate("126-")
	assert.Error(t, valid)

	valid = validate("126--almost--valid")
	assert.Error(t, valid)

	valid = validate("126-almost--valid")
	assert.Error(t, valid)

	valid = validate("126-almost-valid--")
	assert.Error(t, valid)

	valid = validate("15")
	assert.Error(t, valid)

	valid = validate("15-valid")
	assert.NoError(t, valid)

	valid = validate("15-valid-")
	assert.Error(t, valid)

	valid = validate("15-valid-code")
	assert.NoError(t, valid)
}

func TestWindowSizeToDialog(t *testing.T) {
	factor := float32(0.8)
	actual := fyne.NewSize(100, 200)
	expected := fyne.NewSize(actual.Width*factor, actual.Height*factor)
	assert.Equal(t, expected, WindowSizeToDialog(actual))
}
