/*
	processLoop
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/grammar"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// processLoop Process loop block.
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processLoop(tokens vector.Vector, do bool) int {
	contentList := tokens.Sublist(1, len(tokens.Vals)-1)

	// Content is empty?
	if contentList.Vals == nil {
		fract.Error(tokens.Vals[0].(objects.Token), "Content is empty!")
	}

	functionLen := len(i.funcs.Vals)
	_break := false
	kwstate := fract.TypeNone
	iindex := i.index

	//*************
	//    WHILE
	//*************
	if len(contentList.Vals) == 1 ||
		contentList.Vals[1].(objects.Token).Type != fract.TypeIn {
		variableLen := len(i.vars.Vals)

		/* Interpret/skip block. */
	next:
		i.index++
		tokens := i.tokens.Vals[i.index].(vector.Vector)
		condition := i.processCondition(contentList)

		if tokens.Vals[0].(objects.Token).Type == fract.TypeBlockEnd { // Block is ended.
			// Remove temporary variables.
			i.vars.Vals = i.vars.Vals[:variableLen]
			// Remove temporary functions.
			i.funcs.Vals = i.funcs.Vals[:functionLen]

			if _break && condition != grammar.TRUE {
				return kwstate
			}

			i.index = iindex
			goto next
		}

		// Condition is true?
		if condition == grammar.TRUE {
			if do {
				kwstate = i.processTokens(tokens, do)
				if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
					_break = true
					i.skipBlock()
				} else if kwstate == fract.LOOPContinue { // Continue loop?
					i.skipBlock()
				}
			}
		} else {
			_break = true
			i.skipBlock()
			i.index--
		}
		goto next
	}

	//*************
	//   FOREACH
	//*************
	nameToken := contentList.Vals[0].(objects.Token)
	// Name is not name?
	if nameToken.Type != fract.TypeName {
		fract.Error(nameToken, "This is not a valid name!")
	}

	// Name is already defined?
	if i.varIndexByName(nameToken.Value) != -1 {
		fract.Error(nameToken, "Already defined this name!: "+nameToken.Value)
	}

	value := i.processValue(contentList.Sublist(2, len(contentList.Vals)-2))

	// Type is not array?
	if !value.Array {
		fract.Error(contentList.Vals[0].(objects.Token),
			"Foreach loop must defined array value!")
	}

	// Create loop variable.
	variable := objects.Variable{
		Name:  nameToken.Value,
		Const: false,
		Value: objects.Value{
			Array:   false,
			Content: []string{""},
		},
	}

	i.vars.Vals = append(i.vars.Vals, variable)

	variableLen := len(i.vars.Vals)

	if len(value.Content) == 0 {
		i.index++
		i.skipBlock()
	} else {
		for vindex := 0; vindex < len(value.Content); {
			i.index++
			tokens := i.tokens.Vals[i.index].(vector.Vector)

			variable.Value.Content[0] = value.Content[vindex]

			if tokens.Vals[0].(objects.Token).Type == fract.TypeBlockEnd { // Block is ended.
				// Remove temporary variables.
				i.vars.Vals = i.vars.Vals[:variableLen]
				// Remove temporary functions.
				i.funcs.Vals = i.funcs.Vals[:functionLen]

				vindex++
				if _break || vindex == len(value.Content) {
					break
				}
				i.index = iindex

				continue
			}

			// Condition is true?
			if do {
				kwstate = i.processTokens(tokens, do)
				if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
					_break = true
					i.skipBlock()
					i.index--
				} else if kwstate == fract.LOOPContinue { // Continue next?
					i.skipBlock()
					i.index--
				}
			} else {
				_break = true
				i.skipBlock()
				i.index--
			}
		}
	}

	// Remove loop variable.
	i.vars.Vals = i.vars.Vals[:variableLen-1]
	return kwstate
}
