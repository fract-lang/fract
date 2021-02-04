/*
	THE INTERPRETER STRUCT
*/

package interpreter

import (
	"../lexer"
	"../utilities/vector"
)

// Interpreter Interprater of Fract.
type Interpreter struct {
	/* PRIVATE */

	// Parser of this file.
	lexer *lexer.Lexer
	// Variables.
	vars *vector.Vector

	/* PUBLIC */

	/* Type of file. */
	Type int
}
