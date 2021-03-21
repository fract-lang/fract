/*
	processExit Function.
*/

package interpreter

import (
	"os"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// processExit Process exit keyword.
// tokens Tokens to process.
func (i *Interpreter) processExit(tokens []obj.Token) {
	first := tokens[0]

	// Value is not defined?
	if len(tokens) < 2 {
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value),
			"Value is not defined!")
	}

	valueSlice := tokens[1:]
	value := i.processValue(&valueSlice)

	if value.Array {
		fract.Error(first, "Array is not a valid value!")
	} else if value.Type != fract.VALInteger {
		fract.Error(first, "Exit code is only be integer!")
	}

	code, _ := arithmetic.ToInt64(value.Content[0])
	os.Exit(int(code))
}
