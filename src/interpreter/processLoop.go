/*
	processLoop
*/

package interpreter

import (
	"../fract"
	"../fract/name"
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

	_break := false
	iindex := i.index

	// ************
	//    WHILE
	// ************
	if len(contentList.Vals) == 1 || contentList.Vals[1].(objects.Token).Type != fract.TypeIn {
		variableLen := len(i.vars.Vals)

		/* Interpret/skip block. */
		i.index++
		for ; i.index < len(i.tokens.Vals); i.index++ {
			tokens = i.tokens.Vals[i.index].(*vector.Vector)
			condition := i.processCondition(contentList)

			first := tokens.Vals[0].(objects.Token)
			if first.Type == fract.TypeBlockEnd { // Block is ended.
				// Remove temporary variables.
				i.vars.Vals = i.vars.Vals[:variableLen]

				if condition != grammar.TRUE || _break {
					i.subtractBlock(&first)
					return
				}

				i.index = iindex
				continue
			}

			// Condition is true?
			if condition == grammar.TRUE {
				if do {
					kwstate := i.processTokens(tokens, do)
					if kwstate == fract.LOOPBreak { // Break loop?
						_break = true
						i.skipBlock(tokens)
						i.index--
					} else if kwstate == fract.LOOPContinue { // Continue loop?
						i.skipBlock(tokens)
						i.index--
					}
				}
			} else {
				_break = true
				i.skipBlock(tokens)
				i.index--
			}
		}
		return
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
	if name.VarIndexByName(i.vars, nameToken.Value) != -1 {
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
		Type:  grammar.DtInt64,
		Value: []string{""},
	}

	if value.Type == fract.VTFloat {
		variable.Type = grammar.DtFloat64
	}

	i.vars.Vals = append(i.vars.Vals, variable)

	variableLen := len(i.vars.Vals)

	for vindex := 0; vindex < len(value.Content); {
		i.index++
		if i.index >= len(i.tokens.Vals) {
			return
		}
		tokens = i.tokens.Vals[i.index].(*vector.Vector)

		variable.Value[0] = value.Content[vindex]

		first := tokens.Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			// Remove temporary variables.
			i.vars.Vals = i.vars.Vals[:variableLen]

			vindex++
			if _break || vindex == len(value.Content) {
				i.subtractBlock(&first)
				break
			}
			i.index = iindex

			continue
		}
		// Condition is true?
		if do {
			kwstate := i.processTokens(tokens, do)
			if kwstate == fract.LOOPBreak { // Break loop?
				_break = true
				i.skipBlock(tokens)
				i.index--
			} else if kwstate == fract.LOOPContinue { // Continue next?
				i.skipBlock(tokens)
				i.index--
			}
		} else {
			_break = true
			i.skipBlock(tokens)
			i.index--
		}
	}

	// Remove loop variable.
	i.vars.Vals = i.vars.Vals[:variableLen-1]
}
