package util

import (
	"regexp"
	"testing"

	"fyne.io/fyne/v2"

	"github.com/alecthomas/assert/v2"
)

func BenchmarkCodeValidator(b *testing.B) {
	for b.Loop() {
		CodeValidator("125-upset-universe-mistake")
	}
}

var codeValidatorTestcases = []struct {
	in   string
	want error
}{
	{"", nil}, // Do not report invalid when empty.
	{"invalid-code", errInvalidCode},
	{"IOI-", errInvalidCode},
	{"126-", errInvalidCode},
	{"126--almost--valid", errInvalidCode},
	{"126-almost--valid", errInvalidCode},
	{"126-almost-valid-", errInvalidCode},
	{"15", errInvalidCode},
	{"15-valid", nil},
	{"15-valid-", errInvalidCode},
	{"15-valid-code", nil},
	{"15-v", nil},
	{"15-^a", errInvalidCode},

	// Cases found by fuzzing:
	{"-0", errInvalidCode},
	{"0- ", errInvalidCode},
}

func TestCodeValidator(t *testing.T) {
	for _, tc := range codeValidatorTestcases {
		out := CodeValidator(tc.in)
		assert.Equal(t, tc.want, out)
	}
}

func FuzzCodeValidator(f *testing.F) {
	for _, tc := range codeValidatorTestcases {
		f.Add(tc.in)
	}

	isValidCode := regexp.MustCompile(`(^\d+(-[a-zA-Z0-9]+)+$)|(^$)`)
	f.Fuzz(func(t *testing.T, input string) {
		reportedValid := CodeValidator(input) == nil
		appearsCorrect := isValidCode.MatchString(input)

		if reportedValid && !appearsCorrect {
			t.Errorf("Code validator returned no error but input seems invalid")
		} else if !reportedValid && appearsCorrect {
			t.Errorf("Code validator returned error but input seems correct")
		}
	})
}

func TestWindowSizeToDialog(t *testing.T) {
	factor := float32(0.8)
	actual := fyne.NewSize(100, 200)
	expected := fyne.NewSize(actual.Width*factor, actual.Height*factor)
	assert.Equal(t, expected, WindowSizeToDialog(actual))
}
