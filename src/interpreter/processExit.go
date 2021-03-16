/*
	processExit Function.
*/

package interpreter

import (
	"os"

	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/fract/arithmetic"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// processExit Process exit keyword.
// tokens Tokens to process.
func (i *Interpreter) processExit(tokens vector.Vector) {
	first := tokens.Vals[0].(objects.Token)

	// Value is not defined?
	if len(tokens.Vals) < 2 {
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value),
			"Value is not defined!")
	}

	value := i.processValue(&vector.Vector{Vals: tokens.Vals[1:]})

	if value.Array {
		fract.Error(first, "Array is not a valid value!")
	} else if arithmetic.IsFloatValue(value.Content[0]) {
		fract.Error(first, "Exit code is only be integer!")
	}

	code, _ := arithmetic.ToInt64(value.Content[0])
	os.Exit(int(code))
}
