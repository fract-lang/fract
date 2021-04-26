/*
	New Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/internal/lexer"
)

// New Create new instance of Parser.
// path Path of directory.
// fpath Path of destination file.
func New(path, fpath string) *Interpreter {
	return &Interpreter{
		Lexer: &lexer.Lexer{
			File: ReadyFile(fpath),
			Line: 1,
		},
	}
}
