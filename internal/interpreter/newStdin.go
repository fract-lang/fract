package interpreter

import (
	"github.com/fract-lang/fract/internal/lexer"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// NewStdin returns new instance of interpreter from standard input.
func NewStdin(path string) *Interpreter {
	return &Interpreter{
		Lexer: &lexer.Lexer{
			File: &objects.SourceFile{
				Path:  fract.Stdin,
				File:  nil,
				Lines: nil,
			},
			Line: 1,
		},
	}
}
