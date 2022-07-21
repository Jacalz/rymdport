package completion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func generator(match string) []string {
	if match == "" {
		return []string{}
	}

	if byte(match[0]) >= '0' && byte(match[0]) <= '9' {
		return []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	}

	return []string{"zero", "one"}
}

func TestTabCompleter_Next(t *testing.T) {
	cmpl := &TabCompleter{Generate: generator}

	for i := 0; i <= 9; i++ {
		num := string('0' + rune(i))
		next := cmpl.Next(num)
		assert.Equal(t, num, next)
	}

	num := "0"
	next := cmpl.Next(num)
	assert.Equal(t, next, num)

	cmpl.Reset()

	num = "zero"
	next = cmpl.Next(num)
	assert.Equal(t, num, next)

	cmpl.Reset()

	next = cmpl.Next("")
	assert.Empty(t, next)

	num = "none"
	next = cmpl.Next("none")
	assert.Equal(t, num, next)
}

func TestTabCompleter_Previous(t *testing.T) {
	cmpl := &TabCompleter{Generate: generator}

	for i := 9; i >= 0; i-- {
		num := string('0' + rune(i))
		prev := cmpl.Previous(num)
		assert.Equal(t, num, prev)
	}

	num := "9"
	prev := cmpl.Previous(num)
	assert.Equal(t, prev, num)

	cmpl.Reset()

	num = "one"
	prev = cmpl.Previous(num)
	assert.Equal(t, num, prev)

	num = "zero"
	prev = cmpl.Previous(num)
	assert.Equal(t, num, prev)

	num = "one"
	prev = cmpl.Previous(num)
	assert.Equal(t, num, prev)

	cmpl.Reset()

	prev = cmpl.Previous("")
	assert.Empty(t, prev)

	num = "none"
	prev = cmpl.Previous("none")
	assert.Equal(t, num, prev)
}

func TestTabCompleter_Bidirectional(t *testing.T) {
	cmpl := &TabCompleter{Generate: generator}

	num := "0"
	next := cmpl.Next(num)
	assert.Equal(t, next, num)

	num = "1"
	next = cmpl.Next(num)
	assert.Equal(t, next, num)

	num = "0"
	next = cmpl.Previous(num)
	assert.Equal(t, next, num)
}
