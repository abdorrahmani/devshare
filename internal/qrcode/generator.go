package qrcode

import (
	"github.com/fatih/color"
	"github.com/mdp/qrterminal/v3"
	"os"
)

// GenerateQrCodeWithMessage generates a QR code with a custom message.
// It prints the QR code to the console and displays the message in cyan color.
// The QR code contains a URL that starts with "http:://" followed by the provided
func GenerateQrCodeWithMessage(text, message string) {
	color.Cyan(message)
	config := qrterminal.Config{
		HalfBlocks: true,
		Level:      qrterminal.L,
		Writer:     os.Stdout,
	}
	qrterminal.GenerateWithConfig("http:://"+text, config)
}
