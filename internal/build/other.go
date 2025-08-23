//go:build !flatpak

// Package build contains build-time constants.
package build

// IsFlatpak indicates whether the build is for Flatpak.
const IsFlatpak = false
