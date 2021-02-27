/*
	processTokens Function.
*/

package interpreter

import (
	"fmt"

	"../fract"
	"../fract/arithmetic"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// processTokens Process tokens and returns true if block end, returns false if not.
// and returns loop keyword state.
//
// tokens Tokens to process.
// do Do processes?
// nested Is nested?
func (i *Interpreter) processTokens(tokens *vector.Vector, do bool) int {
	first := tokens.Vals[0].(objects.Token)

	if first.Type == fract.TypeBlockEnd {
		i.subtractBlock(&first)
		return -1
	}

	if first.Type == fract.TypeValue || first.Type == fract.TypeBrace ||
		first.Type == fract.TypeName || first.Type == fract.TypeBooleanTrue ||
		first.Type == fract.TypeBooleanFalse {
		// Check variable set statement?
		if first.Type == fract.TypeName {
			for index := range tokens.Vals {
				current := tokens.Vals[index].(objects.Token)
				if current.Type == fract.TypeOperator &&
					current.Value == grammar.Setter { // Variable setting.
					i.processVariableSet(tokens)
					return -1
				}
			}
		}

		// Println
		value := i.processValue(tokens)
		if value.Array {
			if value.Charray {
				for index := range value.Content {
					ch, _ := arithmetic.ToInt64(value.Content[index])
					fmt.Printf("%c", ch)
				}
				fmt.Println()
			} else {
				fmt.Println(value.Content)
			}
		} else {
			if value.Charray {
				ch, _ := arithmetic.ToInt64(value.Content[0])
				fmt.Printf("%c\n", ch)
			} else {
				fmt.Println(value.Content[0])
			}
		}
	} else if first.Type == fract.TypeVariable { // Variable definition.
		i.processVariableDefinition(tokens)
	} else if first.Type == fract.TypeDelete { // Delete from memory.
		i.processDelete(tokens)
	} else if first.Type == fract.TypeIf { // if-elif-else.
		return i.processIf(tokens, do)
	} else if first.Type == fract.TypeLoop { // Loop.
		i.loops++
		i.processLoop(tokens, do)
		i.loops--
	} else if first.Type == fract.TypeBreak { // Break loop.
		if i.loops == 0 {
			fract.Error(first, "Break keyword only used in loops!")
		}
		return fract.LOOPBreak
	} else if first.Type == fract.TypeContinue { // Continue loop.
		if i.loops == 0 {
			fract.Error(first, "Continue keyword only used in loops!")
		}
		return fract.LOOPContinue
	} else if first.Type == fract.TypeExit { // Exit.
		i.processExit(tokens)
	} else if first.Type == fract.TypeFunction { // Function.
		i.processFunction(tokens)
	} else {
		fract.Error(first, "What is this?: "+first.Value)
	}
	return -1
}
