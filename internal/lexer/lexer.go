package lexer

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Lexer of Fract.
type Lexer struct {
	/* PRIVITE */

	// Last generated token.
	lastToken obj.Token
	// Multiline comment process state.
	multilineComment bool
	// Brace count.
	braceCount int
	// Bracket cout.
	bracketCount int
	// Parenthes count.
	parenthesCount int

	/* PUBLIC */

	// Destination file.
	File obj.CodeFile
	// Last column.
	Column int
	// Last line.
	Line int
	// Finished file.
	Finished bool
}
