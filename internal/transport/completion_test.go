package transport

import "testing"

var completions []string

func BenchmarkNameplateCompletion(b *testing.B) {
	c := Client{}

	local := []string{}

	for i := 0; i < b.N; i++ {
		local = c.CompleteRecvCode("5-letterhead-be")
	}

	completions = local
}
