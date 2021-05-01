package interpreter

import (
	"github.com/fract-lang/fract/internal/lexer"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Interpreter of Fract.
type Interpreter struct {
	variables         []*obj.Variable
	functions         []obj.Function
	funcTempVariables int              // Count of function temporary variables.
	loopCount         int
	functionCount     int
	index             int
	returnValue       *obj.Value       // Pointer of last return value.

	Lexer             *lexer.Lexer
	Tokens            [][]obj.Token    // All Tokens of code file.
	Imports           []*ImportInfo
}
