package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed svg/qrcode.svg
var qrcodeData []byte

var QrCodeIcon = &fyne.StaticResource{
	StaticName:    "qrcode.svg",
	StaticContent: qrcodeData,
}
