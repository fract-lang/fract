/*
	processFunction Function.
*/

package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processFunction Process function.
// tokens Tokens to process.
// protected Protected?
func (i *Interpreter) processFunction(tokens []obj.Token, protected bool) {
	tokenLen := len(tokens)
	_name := tokens[1]

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	// Name is already defined?
	if index := i.functionIndexByName(_name.Value); index != -1 {
		fract.Error(_name, "Already defined function in this name at line: "+
			fmt.Sprint(i.functions[index].Line))
	}

	// Function parentheses are not defined?
	if tokenLen < 4 {
		fract.Error(_name, "Where is the function parentheses?")
	}

	i.index++
	function := obj.Function{
		Name:       _name.Value,
		Start:      i.index,
		Line:       _name.Line,
		Parameters: []obj.Parameter{},
		Protected:  protected,
	}

	dtToken := tokens[tokenLen-1]
	if dtToken.Type != fract.TypeBrace ||
		dtToken.Value != grammar.TokenRParenthes {
		fract.Error(dtToken, "Invalid syntax!")
	}

	if paramList := vector.Sublist(tokens, 3, tokenLen-4); paramList != nil {
		paramList := *paramList

		// Decompose function parameters.
		paramName, defaultDefined := true, false
		var lastParameter obj.Parameter
		for index := 0; index < len(paramList); index++ {
			current := paramList[index]
			if paramName {
				if current.Type == fract.TypeParams {
					continue
				}

				if current.Type != fract.TypeName {
					fract.Error(current, "Parameter name is not found!")
				}

				lastParameter = obj.Parameter{
					Name:   current.Value,
					Params: index > 0 && paramList[index-1].Type == fract.TypeParams,
				}
				function.Parameters = append(function.Parameters, lastParameter)
				paramName = false
				continue
			} else {
				paramName = true

				// Default value definition?
				if current.Value == grammar.TokenEquals {
					brace := 0
					index++
					start := index
					for ; index < len(paramList); index++ {
						current = paramList[index]
						if current.Type == fract.TypeBrace {
							if current.Value == grammar.TokenLBrace ||
								current.Value == grammar.TokenLParenthes ||
								current.Value == grammar.TokenLBracket {
								brace++
							} else {
								brace--
							}
						} else if current.Type == fract.TypeComma {
							break
						}
					}
					if index-start < 1 {
						fract.Error(paramList[start-1], "Value is not defined!")
					}
					lastParameter.Default = i.processValue(
						vector.Sublist(paramList, start, index-start))
					if !lastParameter.Default.Array {
						fract.Error(current, "Params parameter is can only take array values!")
					}
					function.Parameters[len(function.Parameters)-1] = lastParameter
					function.DefaultParameterCount++
					defaultDefined = true
					continue
				}

				if lastParameter.Default.Content == nil && defaultDefined {
					fract.Error(current,
						"All parameters after a given parameter with a default value must take a default value!")
				}

				if current.Type != fract.TypeComma {
					fract.Error(current, "Comma is not found!")
				}
			}
		}

		if lastParameter.Default.Content == nil && defaultDefined {
			fract.Error(tokens[len(tokens)-1],
				"All parameters after a given parameter with a default value must take a default value!")
		}
	}

	i.skipBlock(false)
	function.Tokens = i.Tokens[function.Start : function.Start+i.index-function.Start+1]
	i.functions = append(i.functions, function)
}
