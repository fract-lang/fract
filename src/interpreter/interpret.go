/*
	Interpret Function
*/

package interpreter

import (
	"fmt"

	"../fract"
	"../grammar"
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
		if !tokens.Any() {
			continue
		}

		first := tokens.At(0).(objects.Token)

		if first.Type == fract.TypeValue || first.Type == fract.TypeParentheses ||
			first.Type == fract.TypeName || first.Type == fract.TypeBooleanTrue ||
			first.Type == fract.TypeBooleanFalse {
			if first.Type == fract.TypeName && tokens.Len() > 1 {
				second := tokens.At(1).(objects.Token)
				if second.Type == fract.TypeOperator &&
					second.Value == grammar.Setter { // Variable setting.
					i.processVariableSet(&tokens)
					continue
				}
			}

			// Println
			fmt.Println(i.processValue(&tokens).Content)
		} else if first.Type == fract.TypeVariable { // Variable definition.
			i.processVariableDefinition(&tokens)
		} else if first.Type == fract.TypeDelete { // Delete from memory.
			i.processDelete(&tokens)
		} else {
			fract.Error(first, "What is this?: "+first.Value)
		}
	}
}
