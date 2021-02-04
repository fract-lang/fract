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
	ErrorCustom(token.Line, token.Column, message)
}

// ErrorCustom Exit with error.
// line Line of error.
// column Column of error.
// message Message of error.
func ErrorCustom(line int, column int, message string) {
	fmt.Printf("RUNTIME ERROR\nMessage: %s\nLINE: %d\nCOLUMN: %d",
		message, line, column)
	os.Exit(1)
}
