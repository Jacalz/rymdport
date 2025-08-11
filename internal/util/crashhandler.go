package util

import (
	"log"
	"os"
	"runtime/debug"
	"time"
)

// SetUpCrashLogging creates a temporary log where potential panics are logged.
func SetUpCrashLogging() func() {
	f, err := os.CreateTemp("", "rymdport-crash-"+time.Now().Format(time.DateOnly)+"-*.log")
	if err != nil {
		log.Println("On creating temporary crash dump file:", err)
		return func() {}
	}
	defer f.Close()

	if err := debug.SetCrashOutput(f, debug.CrashOptions{}); err != nil {
		log.Println("On setting crash output handler:", err)
	}

	name := f.Name()
	return func() {
		os.Remove(name)
	}
}
