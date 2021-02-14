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
	i.blockCount++
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

	_continue := false
	_break := false

	cacheList := tokens.Sublist(index+1, tokens.Len()-index-1)
	tokens = &cacheList

	i.emptyControl(&tokens)
	iindex := i.index

	// WHILE
	if contentList.Len() == 1 || contentList.At(1).(objects.Token).Type != fract.TypeIn {
		variableLen := i.vars.Len()

		/* Interpret/skip block. */
		for i.index < i.tokenLen {
			i.index++
			tokens = i.tokens.At(i.index).(*vector.Vector)
			condition := i.processCondition(&contentList)

			first := tokens.First().(objects.Token)
			if first.Type == fract.TypeBlockEnd { // Block is ended.
				if condition != grammar.TRUE || _break {
					i.subtractBlock(&first)
					return
				}
				i.index = iindex
				_continue = false

				// Remove temporary variables.
				i.vars.RemoveRange(variableLen, i.vars.Len()-variableLen)

				continue
			}

			// Condition is true?
			if condition == grammar.TRUE {
				if do && !_continue {
					kwstate := i.processTokens(tokens, do)
					if kwstate == fract.LOOPBreak { // Break loop?
						do = false
						_break = true
					} else {
						_continue = kwstate == fract.LOOPContinue // Continue loop?
					}
				}
			} else {
				do = false
				_break = true
				if first.Type == fract.TypeIf { // If?
					i.processIf(tokens, do)
				}
			}
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
		i.index++
		tokens = i.tokens.At(i.index).(*vector.Vector)

		variable.Value[0] = value.Content[vindex]
		i.vars.Set(i.vars.Len()-1, variable)

		first := tokens.First().(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			vindex++
			if _break || vindex == len(value.Content) {
				i.subtractBlock(&first)
				break
			}
			i.index = iindex
			_continue = false

			// Remove temporary variables.
			i.vars.RemoveRange(variableLen, i.vars.Len()-variableLen)

			continue
		}

		// Condition is true?
		if do && !_continue {
			kwstate := i.processTokens(tokens, do)
			if kwstate == fract.LOOPBreak { // Break loop?
				do = false
				_break = true
			} else {
				_continue = kwstate == fract.LOOPContinue // Continue next?
			}
		} else {
			do = false
			_break = true
			if first.Type == fract.TypeIf { // If?
				i.processIf(tokens, do)
			}
		}
	}

	// Remove loop variable.
	i.vars.RemoveLast()
}
