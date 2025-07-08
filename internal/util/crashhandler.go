package util

import (
	"os"
	"runtime/debug"
	"time"
)

// SetUpCrashLogging creates a temporary log where potential panics are logged.
func SetUpCrashLogging() {
	f, err := os.CreateTemp("", "Rymdport "+time.Now().Format(time.DateTime)+".log")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = debug.SetCrashOutput(f, debug.CrashOptions{})
	if err != nil {
		panic(err)
	}
}
