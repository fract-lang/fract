/*
	THE INTERPRETER STRUCT
*/

package interpreter

import (
	"../lexer"
	"../utils/vector"
)

// Interpreter Interprater of Fract.
type Interpreter struct {
	/* PRIVATE */

	// Parser of this file.
	lexer lexer.Lexer
	// Variables.
	vars vector.Vector
	// Functions.
	funcs vector.Vector
	// Count of function temporary variables.
	funcTempVariables int
	// Loop count.
	loops int
	// Function count.
	functions int
	// All tokens of code file.
	tokens vector.Vector
	// Interpreter index.
	index int
	// Last return index.
	returnIndex int

	/* PUBLIC */

	/* Type of file. */
	Type int
}
