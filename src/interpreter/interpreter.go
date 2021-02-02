/*
	THE INTERPRETER STRUCT
*/

package interpreter

import (
	"../lexer"
)

// Interpreter Interprater of Fract.
type Interpreter struct {
	/* PRIVATE */

	// Parser of this file.
	lexer *lexer.Lexer

	/* PUBLIC */

	/* Type of file. */
	Type int
}
