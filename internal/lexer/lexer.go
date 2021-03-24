package lexer

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Lexer of Fract.
type Lexer struct {
	/* PRIVITE */

	// Last generated token.
	lastToken obj.Token

	/* PUBLIC */

	// Destination file.
	File obj.CodeFile
	// Last column.
	Column int
	// Last line.
	Line int
	// Finished file.
	Finished bool
	// RangeComment comment process state.
	RangeComment bool
	// Brace count.
	BraceCount int
	// Bracket cout.
	BracketCount int
	// Parenthes count.
	ParenthesCount int
}
