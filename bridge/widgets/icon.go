package widgets

import (
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

func extractMimeType(uri fyne.URI) string {
	mimeTypeSplit := strings.Split(uri.MimeType(), "/")
	if len(mimeTypeSplit) <= 1 {
		return ""
	}
	return mimeTypeSplit[0]
}

func iconFromURI(uri fyne.URI) fyne.Resource {
	switch extractMimeType(uri) {
	case "application":
		return theme.FileApplicationIcon()
	case "audio":
		return theme.FileAudioIcon()
	case "image":
		return theme.FileImageIcon()
	case "text":
		return theme.FileTextIcon()
	case "video":
		return theme.FileVideoIcon()
	default:
		return theme.DocumentIcon() // Sending text
	}
}
