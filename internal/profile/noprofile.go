//go:build !profile

// Package profile contains tooling for easy profiling.
package profile

func Start() func() {
	return func() {}
}
