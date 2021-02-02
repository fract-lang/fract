/*
	ERROR FUNCTIONS
*/

package fract

import (
	"fmt"
	"os"

	"../objects"
)

// Error Exit with error.
// token Token of error.
// message Message of error.
func Error(token objects.Token, message string) {
	fmt.Printf("RUNTIME ERROR\nMessage: %s\nLINE: %d\nCOLUMN: %d",
		message, token.Line, token.Column)
	os.Exit(1)
}
