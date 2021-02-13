/*
	processLoop
*/

package interpreter

import (
	"../fract"
	"../fract/dt"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processLoop Process loop block.
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processLoop(tokens *vector.Vector, do bool) {
	index := parser.IndexBlockDeclare(tokens)
	// Block declare is not defined?
	if index == -1 {
		fract.Error(tokens.Last().(objects.Token), "Where is the block declare!?")
	}

	contentList := tokens.Sublist(1, index-1)
	// Content is empty?
	if !contentList.Any() {
		fract.Error(tokens.First().(objects.Token), "Content is empty!")
	}

	line := i.lexer.Line
	_continue := false

	cacheList := tokens.Sublist(index+1, tokens.Len()-index-1)
	tokens = &cacheList

	// WHILE
	if contentList.Len() == 1 || contentList.At(1).(objects.Token).Type != fract.TypeIn {
		variableLen := i.vars.Len()

		/* Interpret/skip block. */
		for !i.lexer.Finished {
			// Skip this loop if tokens are empty.
			if !tokens.Any() {
				tokens = i.lexer.Next()
				continue
			}

			first := tokens.First().(objects.Token)
			if first.Type == fract.TypeBlockEnd { // Block is ended.
				if line == -1 {
					return
				}
				_continue = false
				i.lexer.Line -= i.lexer.Line - line
				i.lexer.BlockCount++

				// Remove temporary variables.
				i.vars.RemoveRange(variableLen, i.vars.Len()-variableLen)
				goto nextWhile
			}

			// Condition is true?
			if i.processCondition(&contentList) == grammar.TRUE {
				if do && !_continue {
					kwstate := i.processTokens(tokens, do)
					if kwstate == fract.LOOPBreak { // Break loop?
						line = -1
					} else if kwstate == fract.LOOPContinue { // Continue loop?
						_continue = true
					}
				}
			} else {
				if first.Type == fract.TypeIf { // If?
					i.processIf(tokens, false)
				}
				line = -1
			}

		nextWhile:
			tokens = i.lexer.Next()
		}
	}

	// ************
	//     FOR
	// ************
	nameToken := contentList.First().(objects.Token)
	// Name is not name?
	if nameToken.Type != fract.TypeName {
		fract.Error(nameToken, "This is not a valid name!")
	}

	// Name is already defined?
	if i.checkName(nameToken.Value) {
		fract.Error(nameToken, "Already defined this name!: "+nameToken.Value)
	}

	contentList = contentList.Sublist(2, contentList.Len()-2)
	value := i.processValue(&contentList)

	// Type is not array?
	if !dt.TypeIsArray(value.Type) {
		fract.Error(contentList.First().(objects.Token), "For loop must defined array value!")
	}

	// Create loop variable.
	variable := objects.Variable{
		Name:  nameToken.Value,
		Array: false,
		Const: false,
		Type:  grammar.DtFloat64,
		Value: []string{""},
	}
	i.vars.Append(variable)

	variableLen := i.vars.Len()

	for vindex := 0; vindex < len(value.Content); {
		variable.Value[0] = value.Content[vindex]
		i.vars.Set(i.vars.Len()-1, variable)

		/* Interpret/skip block. */
		for !i.lexer.Finished {
			// Skip this loop if tokens are empty.
			if !tokens.Any() {
				tokens = i.lexer.Next()
				continue
			}
			break
		}

		first := tokens.First().(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			vindex++
			if vindex == len(value.Content) {
				break
			}
			_continue = false
			i.lexer.Line -= i.lexer.Line - line
			i.lexer.BlockCount++

			// Remove temporary variables.
			i.vars.RemoveRange(variableLen, i.vars.Len()-variableLen)

			goto nextFor
		}

		// Condition is true?
		if do && !_continue {
			kwstate := i.processTokens(tokens, do)
			if kwstate == fract.LOOPBreak { // Break loop?
				do = false
			} else if kwstate == fract.LOOPContinue { // Continue next?
				_continue = true
			}
		} else {
			if first.Type == fract.TypeIf { // If?
				i.processIf(tokens, false)
			}
		}

	nextFor:
		tokens = i.lexer.Next()
	}

	// Remove loop variable.
	i.vars.RemoveLast()
}
