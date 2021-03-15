/*
	ErrorCustom Function.
*/

package fract

import (
	"fmt"
	"os"
)

// ErrorCustom Exit with error.
// path File path of error.
// line Line of error.
// column Column of error.
// message Message of error.
func ErrorCustom(path string, line, column int, message string) {
	fmt.Printf("RUNTIME ERROR\nFILE: %s\nMessage: %s\nLINE: %d\nCOLUMN: %d",
		path, message, line, column)
	os.Exit(1)
}
