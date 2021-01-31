package fract

import (
	"fmt"
	"os"

	"../objects"
)

const (
	// FractVersion Version of Fract.
	FractVersion string = "0.0.1"

	// FractExtension File extension of Fract.
	FractExtension string = ".fract"
)

// Error Exit with error.
// token Token of error.
// message Message of error.
func Error(token objects.Token, message string) {
	fmt.Printf("RUNTIME ERROR\nMessage: %s\nLINE: %d\nCOLUMN: %d",
		message, token.Line, token.Column)
	os.Exit(1)
}
