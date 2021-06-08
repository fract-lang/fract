package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

//! A change added here(especially added a code block) must also be compatible with "import.go" and
//! add to "isBlockStatement.go" of parser.

// processTokens returns true if block end, returns false if not and returns keyword state.
func (i *Interpreter) processTokens(tokens []objects.Token) uint8 {
	tokens = append([]objects.Token{}, tokens...)
	switch first := tokens[0]; first.Type {
	case
		fract.TypeValue,
		fract.TypeBrace,
		fract.TypeName:
		if first.Type == fract.TypeName {
			brace := 0
			for _, current := range tokens {
				if current.Type == fract.TypeBrace {
					if current.Value == "{" || current.Value == "[" || current.Value == "(" {
						brace++
					} else {
						brace--
					}
				}
				if brace > 0 {
					continue
				}
				if current.Type == fract.TypeOperator &&
					(current.Value == "=" ||
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
		if value := i.processValue(tokens); fract.InteractiveShell {
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
			value := i.processValue(tokens[1:])
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
	case fract.TypeMacro: // Macro.
		return i.processMacro(tokens)
	case fract.TypeDefer: // Defer.
		if l := len(tokens); l < 2 {
			fract.Error(tokens[0], "Function is not defined!")
		} else if tokens[1].Type != fract.TypeName {
			fract.Error(tokens[1], "Invalid syntax!")
		} else if l < 3 {
			fract.Error(tokens[1], "Invalid syntax!")
		} else if tokens[2].Type != fract.TypeBrace || tokens[2].Value != "(" {
			fract.Error(tokens[2], "Invalid syntax!")
		}
		defers = append(defers, i.processFunctionCallModel(tokens[1:]))
	default:
		fract.Error(first, "Invalid syntax!")
	}
	return fract.TypeNone
}
