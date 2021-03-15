/*
	Error Function.
*/

package fract

import (
	"github.com/fract-lang/fract/src/objects"
)

// Error Exit with error.
// token Token of error.
// message Message of error.
func Error(token objects.Token, message string) {
	ErrorCustom(token.File, token.Line, token.Column, message)
}
