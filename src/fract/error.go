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
	ErrorCustom(token.File.Path, token.Line, token.Column, message)
}

// ErrorCustom Exit with error.
// path File path of error.
// line Line of error.
// column Column of error.
// message Message of error.
func ErrorCustom(path string, line int, column int, message string) {
	fmt.Printf("RUNTIME ERROR\nFILE: %s\nMessage: %s\nLINE: %d\nCOLUMN: %d",
		path, message, line, column)
	os.Exit(1)
}
