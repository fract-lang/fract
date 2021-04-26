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
func (i *Interpreter) processLoop(tokens []obj.Token) int16 {
	// Content is empty?
	if vtokens := vector.Sublist(tokens, 1, len(tokens)-1); vtokens == nil {
		tokens = nil
	} else {
		tokens = *vtokens
	}

	functionLen := len(i.functions)
	_break := false
	kwstate := fract.TypeNone
	iindex := i.index

	// processKwState Process and return return value of kwstate.
	processKwState := func() int16 {
		if kwstate != fract.FUNCReturn {
			return fract.TypeNone
		}
		return kwstate
	}

	//*************
	//    WHILE
	//*************
	if tokens == nil || len(tokens) >= 1 {
		if len(tokens) == 1 || len(tokens) >= 1 && tokens[1].Type != fract.TypeIn && tokens[1].Type != fract.TypeComma {
			variableLen := len(i.variables)

			/* Infinity loop. */
			if tokens == nil {
				for {
					i.index++
					tokens := i.Tokens[i.index]

					if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
						// Remove temporary variables.
						i.variables = i.variables[:variableLen]
						// Remove temporary functions.
						i.functions = i.functions[:functionLen]

						if _break {
							return processKwState()
						}

						i.index = iindex
						continue
					}

					kwstate = i.processTokens(tokens)
					if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
						_break = true
						i.skipBlock(false)
						i.index--
					} else if kwstate == fract.LOOPContinue { // Continue loop?
						i.skipBlock(false)
						i.index--
					}
				}
			}

			/* Interpret/skip block. */
			conditionList := &tokens
			condition := i.processCondition(conditionList)
			for {
				i.index++
				tokens := i.Tokens[i.index]

				if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
					// Remove temporary variables.
					i.variables = i.variables[:variableLen]
					// Remove temporary functions.
					i.functions = i.functions[:functionLen]

					condition = i.processCondition(conditionList)

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
	}

	//*************
	//   FOREACH
	//*************
	nameToken := tokens[0]

	// Name is not name?
	if nameToken.Type != fract.TypeName {
		fract.Error(nameToken, "This is not a valid name!")
	}

	// Element name?
	elementName := ""

	if tokens[1].Type == fract.TypeComma {
		if len(tokens) < 3 || tokens[2].Type != fract.TypeName {
			fract.Error(tokens[1], "Element name is not defined!")
		}

		if tokens[2].Value != grammar.TokenUnderscore {
			elementName = tokens[2].Value
		}

		if len(tokens)-3 == 0 {
			tokens[2].Column += len(tokens[2].Value)
			fract.Error(tokens[2], "Value is not defined!")
		}
		tokens = tokens[2:]
	}

	if vtokens, inToken := vector.Sublist(tokens, 2, len(tokens)-2), tokens[1]; vtokens != nil {
		tokens = *vtokens
	} else {
		fract.Error(inToken, "Value is not defined!")
	}

	value := i.processValue(&tokens)

	// Type is not array?
	if !value.Array && value.Content[0].Type != fract.VALString {
		fract.Error(tokens[0], "Foreach loop must defined array value!")
	}

	// Empty array?
	if len(value.Content) == 0 {
		i.index++
		i.skipBlock(false)
		return kwstate
	}

	i.variables = append(
		[]*obj.Variable{
			{ // Index.
				Name: nameToken.Value,
				Value: obj.Value{
					Content: []obj.DataFrame{{Data: "0"}},
				},
			},
			{ // Element.
				Name:  elementName,
				Value: obj.Value{},
			}}, i.variables...)

	variables := i.variables
	index := i.variables[0]
	element := i.variables[1]

	if index.Name == grammar.TokenUnderscore {
		index.Name = ""
	}

	var length int
	if value.Array {
		length = len(value.Content)
	} else {
		length = len(value.Content[0].Data)
	}

	if element.Name != "" {
		if value.Array {
			element.Value.Content = []obj.DataFrame{value.Content[0]}
		} else {
			element.Value.Content = []obj.DataFrame{{
				Data: string(value.Content[0].Data[0]),
				Type: fract.VALString,
			}}
		}
	}

	//? Interpret block.
	for vindex := 0; vindex < length; {
		i.index++
		tokens := i.Tokens[i.index]

		if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
			// Remove temporary variables.
			i.variables = variables
			// Remove temporary functions.
			i.functions = i.functions[:functionLen]

			vindex++
			if _break ||
				(value.Array && vindex == len(value.Content) ||
					!value.Array && vindex == len(value.Content[0].Data)) {
				break
			}
			i.index = iindex

			if index.Name != "" {
				index.Value.Content[0] = obj.DataFrame{Data: fmt.Sprintf("%d", vindex)}
			}

			if element.Name != "" {
				if value.Array {
					element.Value.Content = []obj.DataFrame{value.Content[vindex]}
				} else {
					element.Value.Content[0] = obj.DataFrame{
						Data: string(value.Content[0].Data[vindex]),
						Type: fract.VALString,
					}
				}
			}
			continue
		}

		kwstate = i.processTokens(tokens)
		if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
			_break = true
			i.skipBlock(false)
			i.index--
		} else if kwstate == fract.LOOPContinue { // Continue next?
			i.skipBlock(false)
			i.index--
		}
	}

	// Remove loop variables.
	i.variables = variables[2:]
	return processKwState()
}
