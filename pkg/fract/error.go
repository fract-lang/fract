package fract

import "github.com/fract-lang/fract/pkg/objects"

// Error thrown exception.
func Error(token objects.Token, message string) {
	ErrorCustom(token.File, token.Line, token.Column, message)
}
