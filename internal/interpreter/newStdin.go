/*
	NewStdin Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/internal/lexer"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// NewStdin Create new instance of interpreter from <stdin>.
// path Path of directory.
func NewStdin(path string) Interpreter {
	return Interpreter{
		Lexer: lexer.Lexer{
			File: obj.CodeFile{
				Path:  fract.Stdin,
				File:  nil,
				Lines: nil,
			},
			Line: 1,
		},
	}
}
