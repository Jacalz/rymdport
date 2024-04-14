package ui

import (
	"os"
	"testing"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var globalTabs *container.AppTabs

func BenchmarkCreate(b *testing.B) {
	b.StopTimer()
	a := app.NewWithID("io.github.jacalz.rymdport")
	w := a.NewWindow("Rymdport")
	var tabs *container.AppTabs
	os.Args = []string{"rymdport"} // Don't read test arguments as uri input.
	b.StartTimer()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tabs = Create(a, w)
	}

	// Don't allow the compiler to optimize out.
	globalTabs = tabs
}
