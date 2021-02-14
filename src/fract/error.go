/*
	Error Function.
*/

package fract

import (
	"../objects"
)

// Error Exit with error.
// token Token of error.
// message Message of error.
func Error(token objects.Token, message string) {
	ErrorCustom(token.File.Path, token.Line, token.Column, message)
}
