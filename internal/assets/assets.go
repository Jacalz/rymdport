package assets

//go:generate fyne bundle -package assets -o bundled.go icon

// AppIcon contains the main application icon.
// TODO: Use go:embed for the next big release.
var AppIcon = resourceIcon512Png
