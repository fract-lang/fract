/*
	Interpret Function
*/

package interpreter

import (
	"fmt"

	"../fract"
	"../objects"
)

// Interpret Interpret file.
func (i *Interpreter) Interpret() {
	// Lexer is finished.
	if i.lexer.Finished {
		return
	}

	/* Interpret all lines. */
	for !i.lexer.Finished {
		tokens := i.lexer.Next()

		// Skip this loop if tokens are empty.
		if tokens.Len() == 0 {
			continue
		}

		first := tokens.At(0).(objects.Token)

		if first.Type == fract.TypeValue || first.Type == fract.TypeBrace ||
			first.Type == fract.TypeName {
			fmt.Println(i.processValue(&tokens).Content)
		} else if first.Type == fract.TypeVariable {
			i.processVariableDefinition(&tokens)
		} else {
			fract.Error(first, "What is this?: "+first.Value)
		}
	}
}
