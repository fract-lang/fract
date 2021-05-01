package interpreter

import (
	"github.com/fract-lang/fract/internal/lexer"
	"github.com/fract-lang/fract/pkg/objects"
)

// Interpreter of Fract.
type Interpreter struct {
	variables         []*objects.Variable
	functions         []objects.Function
	macroDefines      []*objects.Variable
	funcTempVariables int                   // Count of function temporary variables.
	loopCount         int
	functionCount     int
	index             int
	returnValue       *objects.Value        // Pointer of last return value.

	Lexer             *lexer.Lexer
	Tokens            [][]objects.Token     // All Tokens of code file.
	Imports           []*ImportInfo
}
