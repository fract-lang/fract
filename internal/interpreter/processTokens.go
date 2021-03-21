/*
	processTokens Function.
*/

package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// printValue Print value to screen.
// value Value to print.
func printValue(value obj.Value) {
	if value.Array {
		fmt.Print(value.Content)
	} else {
		fmt.Print(value.Content[0])
	}
}

// processTokens Process tokens and returns true if block end, returns false if not.
// and returns loop keyword state.
//
// tokens Tokens to process.
// nested Is nested?
func (i *Interpreter) processTokens(tokens vector.Vector) int {
	tokens = vector.Vector{Vals: append(make([]interface{}, 0), tokens.Vals...)}

	first := tokens.Vals[0].(obj.Token)

	if first.Type == fract.TypeValue || first.Type == fract.TypeBrace ||
		first.Type == fract.TypeName || first.Type == fract.TypeBooleanTrue ||
		first.Type == fract.TypeBooleanFalse {
		if first.Type == fract.TypeName {
			for _, current := range tokens.Vals {
				current := current.(obj.Token)
				if current.Type == fract.TypeBrace {
					break
				} else if current.Type == fract.TypeOperator &&
					(current.Value == grammar.TokenEquals ||
						current.Value == grammar.Input) { // Variable setting.
					i.processVariableSet(tokens)
					return fract.TypeNone
				}
			}
		}

		// Println
		printValue(i.processValue(&tokens))
		fmt.Println()
	} else if first.Type == fract.TypeProtected { // Protected declaration.
		if len(tokens.Vals) < 2 {
			fract.Error(first, "Protected but what is it protected?")
		}
		second := tokens.Vals[1].(obj.Token)
		tokens.Vals = tokens.Vals[1:]
		if second.Type == fract.TypeVariable { // Variable definition.
			i.processVariableDefinition(tokens, true)
		} else if second.Type == fract.TypeFunction { // Function definition.
			i.processFunction(tokens, true)
		} else {
			fract.Error(second, "Syntax error, you can protect only deletable objects!")
		}
	} else if first.Type == fract.TypeVariable { // Variable definition.
		i.processVariableDefinition(tokens, false)
	} else if first.Type == fract.TypeDelete { // Delete from memory.
		i.processDelete(tokens)
	} else if first.Type == fract.TypeIf { // if-elif-else.
		return i.processIf(tokens)
	} else if first.Type == fract.TypeLoop { // Loop definition.
		i.loops++
		state := i.processLoop(tokens)
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
	} else if first.Type == fract.TypeReturn { // Return.
		if i.functions == 0 {
			fract.Error(first, "Return keyword only used in functions!")
		}
		i.returnIndex = i.index
		return fract.FUNCReturn
	} else if first.Type == fract.TypeFunction { // Function definiton.
		i.processFunction(tokens, false)
	} else if first.Type == fract.TypeExit { // Exit.
		i.processExit(tokens)
	} else {
		fract.Error(first, "What is this?")
	}
	return fract.TypeNone
}
