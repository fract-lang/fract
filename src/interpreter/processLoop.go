/*
	processLoop
*/

package interpreter

import (
	"../fract"
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
		fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token), "Where is the block declare!?")
	}

	contentList := tokens.Sublist(1, index-1)
	// Content is empty?
	if len(contentList.Vals) == 0 {
		fract.Error(tokens.Vals[0].(objects.Token), "Content is empty!")
	}

	_continue := false
	_break := false

	tokens = tokens.Sublist(index+1, len(tokens.Vals)-index-1)

	i.emptyControl(&tokens)
	iindex := i.index

	// WHILE
	if len(contentList.Vals) == 1 || contentList.Vals[1].(objects.Token).Type != fract.TypeIn {
		variableLen := len(i.vars.Vals)

		/* Interpret/skip block. */
		for i.index < len(i.tokens.Vals) {
			i.index++
			tokens = i.tokens.Vals[i.index].(*vector.Vector)
			condition := i.processCondition(contentList)

			first := tokens.Vals[0].(objects.Token)
			if first.Type == fract.TypeBlockEnd { // Block is ended.
				if condition != grammar.TRUE || _break {
					i.subtractBlock(&first)
					return
				}
				i.index = iindex
				_continue = false

				// Remove temporary variables.
				i.vars.Vals = i.vars.Vals[:variableLen]

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
	nameToken := contentList.Vals[0].(objects.Token)
	// Name is not name?
	if nameToken.Type != fract.TypeName {
		fract.Error(nameToken, "This is not a valid name!")
	}

	// Name is already defined?
	if i.checkName(nameToken.Value) {
		fract.Error(nameToken, "Already defined this name!: "+nameToken.Value)
	}

	value := i.processValue(contentList.Sublist(2, len(contentList.Vals)-2))

	// Type is not array?
	if !value.Array {
		fract.Error(contentList.Vals[0].(objects.Token), "For loop must defined array value!")
	}
	// Create loop variable.
	variable := objects.Variable{
		Name:  nameToken.Value,
		Array: false,
		Const: false,
		Type:  grammar.DtFloat64,
		Value: []string{""},
	}
	i.vars.Vals = append(i.vars.Vals, variable)

	variableLen := len(i.vars.Vals)

	for vindex := 0; vindex < len(value.Content); {
		i.index++
		tokens = i.tokens.Vals[i.index].(*vector.Vector)

		variable.Value[0] = value.Content[vindex]
		i.vars.Vals[len(i.vars.Vals)-1] = variable

		first := tokens.Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			vindex++
			if _break || vindex == len(value.Content) {
				i.subtractBlock(&first)
				break
			}
			i.index = iindex
			_continue = false

			// Remove temporary variables.
			i.vars.Vals = i.vars.Vals[:variableLen]

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
	i.vars.Vals = i.vars.Vals[:len(i.vars.Vals)-1]
}
