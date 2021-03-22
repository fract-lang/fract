/*
	processLoop
*/

package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processLoop Process loop block.
// tokens Tokens to process.
func (i *Interpreter) processLoop(tokens []obj.Token) int {
	contentList := *vector.Sublist(tokens, 1, len(tokens)-1)

	// Content is empty?
	if contentList == nil {
		fract.Error(tokens[0], "Content is empty!")
	}

	functionLen := len(i.funcs)
	_break := false
	kwstate := fract.TypeNone
	iindex := i.index

	// processKwState Process and return return value of kwstate.
	processKwState := func() int {
		if kwstate != fract.FUNCReturn {
			return fract.TypeNone
		}
		return kwstate
	}

	//*************
	//    WHILE
	//*************
	if len(contentList) == 1 ||
		contentList[1].Type != fract.TypeIn {
		variableLen := len(i.vars)

		/* Interpret/skip block. */
		for {
			i.index++
			tokens := i.tokens[i.index]
			condition := i.processCondition(&contentList)

			if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
				// Remove temporary variables.
				i.vars = i.vars[:variableLen]
				// Remove temporary functions.
				i.funcs = i.funcs[:functionLen]

				if _break || condition != grammar.KwTrue {
					return processKwState()
				}

				i.index = iindex
				continue
			}

			// Condition is true?
			if condition == grammar.KwTrue {
				kwstate = i.processTokens(tokens)
				if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
					_break = true
					i.skipBlock(false)
					i.index--
				} else if kwstate == fract.LOOPContinue { // Continue loop?
					i.skipBlock(false)
					i.index--
				}
			} else {
				_break = true
				i.skipBlock(false)
				i.index--
			}
		}
	}

	//*************
	//   FOREACH
	//*************
	nameToken := contentList[0]
	// Name is not name?
	if nameToken.Type != fract.TypeName {
		fract.Error(nameToken, "This is not a valid name!")
	}

	// Name is already defined?
	if index := i.varIndexByName(nameToken.Value); index != -1 {
		fract.Error(nameToken, "Already defined variable in this name at line: "+
			fmt.Sprint(i.vars[index].Line))
	}

	inToken := contentList[1]
	contentList = *vector.Sublist(contentList, 2, len(contentList)-2)

	// Value is not defined?
	if contentList == nil {
		fract.Error(inToken, "Value is not defined!")
	}

	value := i.processValue(&contentList)

	// Type is not array?
	if !value.Array {
		fract.Error(contentList[2], "Foreach loop must defined array value!")
	}

	// Empty array?
	if len(value.Content) == 0 {
		i.index++
		i.skipBlock(false)
		return kwstate
	}

	// Create loop variable.
	variable := obj.Variable{
		Name:  nameToken.Value,
		Const: false,
		Value: obj.Value{
			Array:   false,
			Content: []string{""},
		},
	}

	i.vars = append(i.vars, variable)

	variableLen := len(i.vars)

	for vindex := 0; vindex < len(value.Content); {
		i.index++
		tokens := i.tokens[i.index]

		variable.Value.Content[0] = value.Content[vindex]

		if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
			// Remove temporary variables.
			i.vars = i.vars[:variableLen]
			// Remove temporary functions.
			i.funcs = i.funcs[:functionLen]

			vindex++
			if _break || vindex == len(value.Content) {
				break
			}
			i.index = iindex

			continue
		}

		// Condition is true?
		kwstate = i.processTokens(tokens)
		if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
			_break = true
			i.skipBlock(false)
			i.index--
			continue
		} else if kwstate == fract.LOOPContinue { // Continue next?
			i.skipBlock(false)
			i.index--
		}
	}

	// Remove loop variable.
	i.vars = i.vars[:variableLen-1]
	return processKwState()
}
