/*
	Error Function.
*/

package fract

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Error Exit with error.
// token Token of error.
// message Message of error.
func Error(token obj.Token, message string) {
	ErrorCustom(token.File, token.Line, token.Column, message)
}
