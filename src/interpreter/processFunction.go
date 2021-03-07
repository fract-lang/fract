/*
	processFunction Function.
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

// processFunction Process function.
// tokens Tokens to process.
func (i *Interpreter) processFunction(tokens *vector.Vector) {
	index := parser.IndexBlockDeclare(tokens)
	// Block declare is not defined?
	if index == -1 {
		fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token),
			"Where is the block declare!?")
	}

	tokenLen := len(tokens.Vals)

	_name := tokens.Vals[1].(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	// Name is already defined?
	if name.FunctionIndexByName(i.funcs, _name.Value) != -1 {
		fract.Error(_name, "Already defined this function!: "+_name.Value)
	}

	// Function parentheses are not defined?
	if tokenLen < 4 {
		fract.Error(_name, "Where is the function parentheses?")
	}

	i.index++
	function := objects.Function{
		Name:       _name.Value,
		Start:      i.index,
		Parameters: []string{},
	}

	dtToken := tokens.Vals[index-1].(objects.Token)
	if dtToken.Type != fract.TypeBrace ||
		dtToken.Value != grammar.TokenRParenthes {
		fract.Error(dtToken, "Invalid syntax!")
	}

	paramList := tokens.Sublist(3, index-4)

	// Decompose function parameters.
	paramName := true
	for index := range paramList.Vals {
		current := paramList.Vals[index].(objects.Token)
		if paramName {
			if current.Type != fract.TypeName {
				fract.Error(current, "Parameter name is not found!")
			}
			function.Parameters = append(function.Parameters,
				paramList.Vals[index].(objects.Token).Value)
			paramName = false
		} else {
			if current.Type != fract.TypeComma {
				fract.Error(current, "Comma is not found!")
			}
			paramName = true
		}
	}

	i.skipBlock()
	function.Tokens = i.tokens.Sublist(function.Start, i.index-function.Start)
	i.funcs.Vals = append(i.funcs.Vals, function)

}
