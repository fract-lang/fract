/*
	New Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/internal/lexer"
)

// New Create new instance of Parser.
// path Path of destination file.
// type Type of file.
func New(path string, _type int) Interpreter {
	return Interpreter{
		lexer: lexer.Lexer{
			File: ReadyFile(path),
			Line: 1,
		},
		Type: _type,
	}
}
