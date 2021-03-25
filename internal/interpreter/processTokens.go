/*
	processTokens Function.
*/

package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// processTokens Process tokens and returns true if block end, returns false if not.
// and returns loop keyword state.
//
// tokens Tokens to process.
// nested Is nested?
func (i *Interpreter) processTokens(tokens []obj.Token) int {
	tokens = append(make([]obj.Token, 0), tokens...)

	first := tokens[0]

	if first.Type == fract.TypeValue || first.Type == fract.TypeBrace ||
		first.Type == fract.TypeName || first.Type == fract.TypeBooleanTrue ||
		first.Type == fract.TypeBooleanFalse {
		if first.Type == fract.TypeName {
			for _, current := range tokens {
				if current.Type == fract.TypeBrace {
					break
				} else if current.Type == fract.TypeOperator &&
					current.Value == grammar.TokenEquals { // Variable setting.
					i.processVariableSet(tokens)
					return fract.TypeNone
				}
			}
		}

		// Println
		if fract.PrintValue(i.processValue(&tokens)) { // Printed?
			fmt.Println()
		}
	} else if first.Type == fract.TypeProtected { // Protected declaration.
		if len(tokens) < 2 {
			fract.Error(first, "Protected but what is it protected?")
		}
		second := tokens[1]
		tokens = tokens[1:]
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
		i.loopCount++
		state := i.processLoop(tokens)
		i.loopCount--
		return state
	} else if first.Type == fract.TypeBreak { // Break loop.
		if i.loopCount == 0 {
			fract.Error(first, "Break keyword only used in loops!")
		}
		return fract.LOOPBreak
	} else if first.Type == fract.TypeContinue { // Continue loop.
		if i.loopCount == 0 {
			fract.Error(first, "Continue keyword only used in loops!")
		}
		return fract.LOOPContinue
	} else if first.Type == fract.TypeReturn { // Return.
		if i.functionCount == 0 {
			fract.Error(first, "Return keyword only used in functions!")
		}
		i.returnIndex = i.index
		return fract.FUNCReturn
	} else if first.Type == fract.TypeFunction { // Function definiton.
		i.processFunction(tokens, false)
	} else if first.Type == fract.TypeTry { // Try-Catch.
		return i.processTryCatch(tokens)
	} else {
		fract.Error(first, "Invalid syntax!")
	}
	return fract.TypeNone
}
