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

	// Variables.
	variables []obj.Variable
	// Functions.
	functions []obj.Function
	// Count of function temporary variables.
	funcTempVariables int
	// Loop count.
	loopCount int
	// Function count.
	functionCount int
	// Interpreter index.
	index int
	// Pointer of last return value.
	returnValue *obj.Value

	/* PUBLIC */

	// Parser of this file.
	Lexer lexer.Lexer
	// All Tokens of code file.
	Tokens [][]obj.Token
	// All imported script files.
	Imports []ImportInfo
}
