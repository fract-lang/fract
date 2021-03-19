/*
	Error Function.
*/

package fract

import (
	"github.com/fract-lang/fract/pkg/objects"
)

// Error Exit with error.
// token Token of error.
// message Message of error.
func Error(token objects.Token, message string) {
	ErrorCustom(token.File, token.Line, token.Column, message)
}
