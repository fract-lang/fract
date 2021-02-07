/*
	processTokens Functions.
*/

package interpreter

import (
	"fmt"

	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// processTokens Process tokens and returns true if block end, returns false if not.
// tokens Tokens to process.
func (i *Interpreter) processTokens(tokens *vector.Vector) {
	// Skip this loop if tokens are empty.v
	if !tokens.Any() {
		return
	}
	first := tokens.At(0).(objects.Token)

	if first.Type == fract.TypeValue || first.Type == fract.TypeBrace ||
		first.Type == fract.TypeName || first.Type == fract.TypeBooleanTrue ||
		first.Type == fract.TypeBooleanFalse {
		if first.Type == fract.TypeName && tokens.Len() > 1 {
			second := tokens.At(1).(objects.Token)
			if second.Type == fract.TypeOperator &&
				second.Value == grammar.Setter { // Variable setting.
				i.processVariableSet(tokens)
				return
			}
		}

		// Println
		fmt.Println(i.processValue(tokens).Content)
	} else if first.Type == fract.TypeVariable { // Variable definition.
		i.processVariableDefinition(tokens)
	} else if first.Type == fract.TypeDelete { // Delete from memory.
		i.processDelete(tokens)
	} else if first.Type == fract.TypeIf { // if-elif-else.
		i.processIf(tokens)
	} else {
		fract.Error(first, "What is this?: "+first.Value)
	}
}
