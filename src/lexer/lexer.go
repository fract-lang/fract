/*
	THE LEXER STRUCT
*/

package lexer

import (
	"../objects"
)

// Lexer of Fract.
type Lexer struct {
	/* PRIVITE */

	// Last generated token.
	lastToken objects.Token

	// Bracket count.
	braceCount int

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
