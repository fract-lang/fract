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

// printValue Print value to screen.
// value Value to print.
func printValue(value objects.Value) {
	if value.Content == nil {
		return
	}

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
}

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
		return fract.TypeNone
	} else if first.Type == fract.TypeValue || first.Type == fract.TypeBrace ||
		first.Type == fract.TypeName || first.Type == fract.TypeBooleanTrue ||
		first.Type == fract.TypeBooleanFalse {
		tokenLen := len(tokens.Vals)
		if tokenLen > 1 {
			second := tokens.Vals[1].(objects.Token)
			// Check name statement?
			if first.Type == fract.TypeName {
				if second.Type == fract.TypeOperator &&
					second.Value == grammar.Setter { // Variable setting.
					i.processVariableSet(tokens)
					return -1
				} else if second.Type == fract.TypeBrace &&
					second.Value == grammar.TokenLParenthes { // Function call.
					i.functions++
					printValue(i.processFunctionCall(tokens))
					i.functions--
					return fract.TypeNone
				}
			}
		}

		// Println
		printValue(i.processValue(tokens))
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
	} else if first.Type == fract.TypeReturn { // Return.
		if i.functions == 0 {
			fract.Error(first, "Return keyword only used in functions!")
		}
		return fract.FUNCReturn
	} else {
		fract.Error(first, "What is this?: "+first.Value)
	}
	return fract.TypeNone
}
