package interpreter

import "github.com/fract-lang/fract/internal/lexer"

// New returns instance of interpreter related to file.
func New(path, fpath string) *Interpreter {
	return &Interpreter{
		Lexer: &lexer.Lexer{
			File: ReadyFile(fpath),
			Line: 1,
		},
	}
}
