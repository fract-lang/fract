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

	switch first := tokens[0]; first.Type {
	case
		fract.TypeValue,
		fract.TypeBrace,
		fract.TypeName,
		fract.TypeBooleanTrue,
		fract.TypeBooleanFalse:
		if first.Type == fract.TypeName {
			brace := 0
			for _, current := range tokens {
				if current.Type == fract.TypeBrace {
					if current.Value == grammar.TokenLBrace ||
						current.Value == grammar.TokenLBracket ||
						current.Value == grammar.TokenLParenthes {
						brace++
					} else {
						brace--
					}
				}
				if brace > 0 {
					continue
				}
				if current.Type == fract.TypeOperator &&
					(current.Value == grammar.TokenEquals ||
						current.Value == grammar.AdditionAssignment ||
						current.Value == grammar.SubtractionAssignment ||
						current.Value == grammar.MultiplicationAssignment ||
						current.Value == grammar.DivisionAssignment ||
						current.Value == grammar.ModulusAssignment ||
						current.Value == grammar.XOrAssignment ||
						current.Value == grammar.LeftBinaryShiftAssignment ||
						current.Value == grammar.RightBinaryShiftAssignment ||
						current.Value == grammar.InclusiveOrAssignment ||
						current.Value == grammar.AndAssignment) { // Variable setting.
					i.processVariableSet(tokens)
					return fract.TypeNone
				}
			}
		}

		// Print value if live interpreting.
		if value := i.processValue(&tokens); fract.LiveInterpret {
			if fract.PrintValue(value) {
				fmt.Println()
			}
		}
	case fract.TypeProtected: // Protected declaration.
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
	case fract.TypeVariable: // Variable definition.
		i.processVariableDefinition(tokens, false)
	case fract.TypeDelete: // Delete from memory.
		i.processDelete(tokens)
	case fract.TypeIf: // if-elif-else.
		return i.processIf(tokens)
	case fract.TypeLoop: // Loop definition.
		i.loopCount++
		state := i.processLoop(tokens)
		i.loopCount--
		return state
	case fract.TypeBreak: // Break loop.
		if i.loopCount == 0 {
			fract.Error(first, "Break keyword only used in loops!")
		}
		return fract.LOOPBreak
	case fract.TypeContinue: // Continue loop.
		if i.loopCount == 0 {
			fract.Error(first, "Continue keyword only used in loops!")
		}
		return fract.LOOPContinue
	case fract.TypeReturn: // Return.
		if i.functionCount == 0 {
			fract.Error(first, "Return keyword only used in functions!")
		}

		if len(tokens) > 1 {
			valueList := tokens[1:]
			value := i.processValue(&valueList)
			i.returnValue = &value
		} else {
			i.returnValue = nil
		}

		return fract.FUNCReturn
	case fract.TypeFunction: // Function definiton.
		i.processFunction(tokens, false)
	case fract.TypeTry: // Try-Catch.
		return i.processTryCatch(tokens)
	case fract.TypeImport: // Import.
		i.processImport(tokens)
	default:
		fract.Error(first, "Invalid syntax!")
	}

	return fract.TypeNone
}
