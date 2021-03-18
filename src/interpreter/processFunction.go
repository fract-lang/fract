/*
	processFunction Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/grammar"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// processFunction Process function.
// tokens Tokens to process.
func (i *Interpreter) processFunction(tokens vector.Vector) {
	tokenLen := len(tokens.Vals)
	_name := tokens.Vals[1].(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	// Name is already defined?
	if i.functionIndexByName(_name.Value) != -1 {
		fract.Error(_name, "Already defined function in this name!: "+_name.Value)
	}

	// Function parentheses are not defined?
	if tokenLen < 4 {
		fract.Error(_name, "Where is the function parentheses?")
	}

	i.index++
	function := objects.Function{
		Name:       _name.Value,
		Start:      i.index,
		Parameters: []objects.Parameter{},
	}

	dtToken := tokens.Vals[tokenLen-1].(objects.Token)
	if dtToken.Type != fract.TypeBrace ||
		dtToken.Value != grammar.TokenRParenthes {
		fract.Error(dtToken, "Invalid syntax!")
	}

	paramList := tokens.Sublist(3, tokenLen-4)

	// Decompose function parameters.
	paramName, defaultDefined := true, false
	var lastParameter objects.Parameter
	for index := 0; index < len(paramList.Vals); index++ {
		current := paramList.Vals[index].(objects.Token)
		if paramName {
			if current.Type != fract.TypeName {
				fract.Error(current, "Parameter name is not found!")
			}

			lastParameter = objects.Parameter{Name: current.Value}
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
				for ; index < len(paramList.Vals); index++ {
					current = paramList.Vals[index].(objects.Token)
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
					fract.Error(paramList.Vals[start-1].(objects.Token),
						"Value is not defined!")
				}
				lastParameter.Default = i.processValue(
					paramList.Sublist(start, index-start))
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
		fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token),
			"All parameters after a given parameter with a default value must take a default value!")
	}

	i.skipBlock(false)
	function.Tokens = i.tokens.Sublist(function.Start, i.index-function.Start+1).Vals
	i.funcs.Vals = append(i.funcs.Vals, function)
}
