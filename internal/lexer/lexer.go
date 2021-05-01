package lexer

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Lexer of Fract.
type Lexer struct {
	lastToken      obj.Token

	File           *obj.CodeFile
	Column         int            // Last column.
	Line           int            // Last line.
	Finished       bool
	RangeComment   bool
	BraceCount     int
	BracketCount   int
	ParenthesCount int
}
