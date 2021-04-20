/*
	Error Function.
*/

package fract

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Error Create new exception.
// token Token of error.
// message Message of error.
func Error(token obj.Token, message string) obj.Exception {
	return ErrorCustom(token.File, token.Line, token.Column, message)
}
