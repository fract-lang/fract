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
	// Loop count.
	loops int

	/* PUBLIC */

	/* Type of file. */
	Type int
}
