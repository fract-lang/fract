/*
	processDelete Function.
*/

package interpreter

import (
	"../fract"
	"../fract/name"
	"../objects"
	"../utilities/vector"
)

// processDelete Process delete keyword.
// tokens Tokens to process.
func (i *Interpreter) processDelete(tokens *vector.Vector) {
	// Value is not defined?
	if tokens.Len() < 2 {
		del := tokens.First().(objects.Token)
		fract.ErrorCustom(del.File.Path, del.Line, del.Column+len(del.Value),
			"Value is not found!")
	}

	for index := 1; index < tokens.Len(); index++ {
		current := tokens.At(index).(objects.Token)

		// Token is not a deletable object?
		if current.Type != fract.TypeName {
			fract.Error(current, "This is not deletable object!")
		}

		position := name.VarIndexByName(i.vars, current.Value)

		// Name is not defined?
		if position == -1 {
			fract.Error(current, "Name is not defined!")
		}

		i.vars.RemoveRange(position, 1)
	}
}
