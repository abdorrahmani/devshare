package qrcode

import (
	"github.com/fatih/color"
	"github.com/mdp/qrterminal/v3"
	"os"
)

func GenerateQrCode(text string) {
	GenerateQrCodeWithMessage(text, "📱 Scan this on your phone:")
}

func GenerateQrCodeWithMessage(text, message string) {
	color.Cyan(message)
	config := qrterminal.Config{
		HalfBlocks: true,
		Level:      qrterminal.L,
		Writer:     os.Stdout,
	}
	qrterminal.GenerateWithConfig(text, config)
}
