/*
	THE INTERPRETER STRUCT
*/

package interpreter

import (
	"github.com/fract-lang/fract/internal/lexer"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Interpreter Interprater of Fract.
type Interpreter struct {
	/* PRIVATE */

	// Parser of this file.
	lexer lexer.Lexer
	// Variables.
	vars []obj.Variable
	// Functions.
	funcs []obj.Function
	// Count of function temporary variables.
	funcTempVariables int
	// Loop count.
	loops int
	// Function count.
	functions int
	// All tokens of code file.
	tokens [][]obj.Token
	// Interpreter index.
	index int
	// Last return index.
	returnIndex int

	/* PUBLIC */

	/* Type of file. */
	Type int
}
