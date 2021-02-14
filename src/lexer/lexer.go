package lexer

import (
	"../objects"
)

// Lexer of Fract.
type Lexer struct {
	/* PRIVITE */

	// Last generated token.
	lastToken objects.Token
	// Brace count.
	braceCount int
	// Bracket cout.
	bracketCount int
	// Parenthes count.
	parenthesCount int

	/* PUBLIC */

	// Destination file.
	File *objects.CodeFile
	// Last column.
	Column int
	// Last line.
	Line int
	// Finished file.
	Finished bool
}
