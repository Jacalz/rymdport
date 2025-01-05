package transport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var completions []string

func BenchmarkNameplateCompletion(b *testing.B) {
	c := Client{}

	local := []string{}

	for range b.N {
		local = c.GenerateCodeCompletion("1-letterhead-be")
	}

	completions = local
}

func TestCompletionGeneration_Progressive(t *testing.T) {
	c := Client{}

	expected := []string{"1-uncut", "1-unearth", "1-unwind", "1-uproot", "1-upset", "1-upshot"}
	assert.Equal(t, expected, c.GenerateCodeCompletion("1-u"))

	expected = []string{"1-uproot", "1-upset", "1-upshot"}
	assert.Equal(t, expected, c.GenerateCodeCompletion("1-up"))

	expected = []string{"1-upset", "1-upshot"}
	assert.Equal(t, expected, c.GenerateCodeCompletion("1-ups"))

	expected = []string{"1-upset"}
	assert.Equal(t, expected, c.GenerateCodeCompletion("1-upse"))

	expected = []string{"1-upshot"}
	assert.Equal(t, expected, c.GenerateCodeCompletion("1-upsh"))

	expected = []string{"1-upset-unicorn", "1-upset-unify", "1-upset-universe"}
	assert.Equal(t, expected, c.GenerateCodeCompletion("1-upset-uni"))
}
