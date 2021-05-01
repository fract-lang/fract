package lexer

import "github.com/fract-lang/fract/pkg/objects"

// Lexer of Fract.
type Lexer struct {
	lastToken      objects.Token

	File           *objects.SourceFile
	Column         int  					     // Last column.
	Line           int            // Last line.
	Finished       bool
	RangeComment   bool
	BraceCount     int
	BracketCount   int
	ParenthesCount int
}
