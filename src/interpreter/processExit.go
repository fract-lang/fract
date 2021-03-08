/*
	processExit Function.
*/

package interpreter

import (
	"os"

	"../fract"
	"../fract/arithmetic"
	"../objects"
	"../utilities/vector"
)

// processExit Process exit keyword.
// tokens Tokens to process.
func (i *Interpreter) processExit(tokens *vector.Vector) {
	first := tokens.Vals[0].(objects.Token)

	// Value is not defined?
	if len(tokens.Vals) < 2 {
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Value is not found!")
	}

	value := i.processValue(vector.New(tokens.Vals[1:]...))

	if value.Array {
		fract.Error(first, "Array is not a valid value!")
	} else if arithmetic.IsFloatValue(value.Content[0]) {
		fract.Error(first, "Exit code is only be integer!")
	}

	code, _ := arithmetic.ToInt64(value.Content[0])
	os.Exit(int(code))
}
