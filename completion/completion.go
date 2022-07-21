// Package completion contains code for handling string completion.
package completion

// TabCompleter handles stepping through a set of completion suggestions.
type TabCompleter struct {
	Generate func(string) []string
	active   []string
	index    int
}

// Next steps into the next completion match and wraps around if necessary.
func (t *TabCompleter) Next(match string) string {
	if t.active == nil && t.Generate != nil {
		t.active = t.Generate(match)
		if len(t.active) > 0 {
			return t.active[0]
		}
	}

	if len(t.active) == 0 {
		return match
	}

	t.index++
	if t.index == len(t.active) {
		t.index = 0
	}

	return t.active[t.index]
}

// Previous steps back into the previous completion match and wraps around if necessary.
func (t *TabCompleter) Previous(match string) string {
	if t.active == nil && t.Generate != nil {
		t.active = t.Generate(match)
		if len(t.active) > 0 {
			t.index = len(t.active) - 1
			return t.active[t.index]
		}
	}

	if len(t.active) == 0 {
		return match
	}

	t.index--
	if t.index < 0 {
		t.index = len(t.active) - 1
	}

	return t.active[t.index]
}

// Reset resets the completion to generate a new on the next call to Next().
func (t *TabCompleter) Reset() {
	t.active = nil
	t.index = 0
}
