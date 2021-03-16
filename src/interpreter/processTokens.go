/*
	processTokens Function.
*/

package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/fract/arithmetic"
	"github.com/fract-lang/fract/src/grammar"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// printValue Print value to screen.
// value Value to print.
func printValue(value objects.Value) bool {
	if value.Type != fract.VALString && value.Content == nil {
		return false
	}

	if value.Array {
		if value.Type == fract.VALString {
			for _, current := range value.Content {
				ch, _ := arithmetic.ToInt64(current)
				fmt.Printf("%c", ch)
			}
		} else {
			if value.Type == fract.VALBoolean {
				if value.Content[0] == "1" {
					fmt.Print(grammar.KwTrue)
				} else {
					fmt.Print(grammar.KwFalse)
				}
			} else {
				fmt.Print(value.Content)
			}
		}
	} else {
		if value.Type == fract.VALString {
			ch, _ := arithmetic.ToInt64(value.Content[0])
			fmt.Printf("%c\n", ch)
		} else {
			fmt.Print(value.Content[0])
		}
	}
	return true
}

// processTokens Process tokens and returns true if block end, returns false if not.
// and returns loop keyword state.
//
// tokens Tokens to process.
// do Do processes?
// nested Is nested?
func (i *Interpreter) processTokens(tokens vector.Vector, do bool) int {
	tokens = vector.Vector{Vals: append(make([]interface{}, 0), tokens.Vals...)}

	first := tokens.Vals[0].(objects.Token)

	if first.Type == fract.TypeBlockEnd {
		fract.Error(first, "The extra block end defined!")
		return fract.TypeNone
	} else if first.Type == fract.TypeValue || first.Type == fract.TypeBrace ||
		first.Type == fract.TypeName || first.Type == fract.TypeBooleanTrue ||
		first.Type == fract.TypeBooleanFalse {
		if first.Type == fract.TypeName {
			for _, current := range tokens.Vals {
				current := current.(objects.Token)
				if current.Type == fract.TypeOperator &&
					(current.Value == grammar.TokenEquals ||
						current.Value == grammar.Input) { // Variable setting.
					i.processVariableSet(tokens)
					return -1
				}
			}
		}

		// Println
		if printValue(i.processValue(&tokens)) { // If printed?
			fmt.Println()
		}
	} else if first.Type == fract.TypeVariable { // Variable definition.
		i.processVariableDefinition(tokens)
	} else if first.Type == fract.TypeDelete { // Delete from memory.
		i.processDelete(tokens)
	} else if first.Type == fract.TypeIf { // if-elif-else.
		return i.processIf(tokens, do)
	} else if first.Type == fract.TypeLoop { // Loop.
		i.loops++
		state := i.processLoop(tokens, do)
		i.loops--
		return state
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
		i.returnIndex = i.index
		return fract.FUNCReturn
	} else {
		fract.Error(first, "What is this?")
	}
	return fract.TypeNone
}
