package ui

import (
	"os"
	"testing"

	"fyne.io/fyne/v2/app"
)

func BenchmarkCreate(b *testing.B) {
	a := app.NewWithID("io.github.jacalz.rymdport")
	w := a.NewWindow("Rymdport")
	os.Args = []string{"rymdport"} // Don't read test arguments as uri input.

	b.ReportAllocs()
	for b.Loop() {
		Create(a, w)
	}
}
