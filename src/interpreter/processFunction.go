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
		fract.Error(_name, "Already defined this name!: "+_name.Value)
	}

	// Function parentheses are not defined?
	if tokenLen < 4 {
		fract.Error(_name, "Where is the function parentheses?")
	}

	i.index++
	function := objects.Function{
		Name:       _name.Value,
		Start:      i.index,
		Parameters: vector.New(),
	}

	dtToken := tokens.Vals[index-1].(objects.Token)
	if dtToken.Type == fract.TypeDataType { // Returnable function?
		index--
		function.ReturnType = dtToken.Value
	} else if dtToken.Type != fract.TypeBrace ||
		dtToken.Value != grammar.TokenRParenthes {
		fract.Error(dtToken, "Invalid syntax!")
	}

	paramList := tokens.Sublist(3, index-4)

	// Decompose function parameters.
	paramName := true
	paramType := true
	for index := range paramList.Vals {
		current := paramList.Vals[index].(objects.Token)
		if paramName {
			if current.Type != fract.TypeName {
				fract.Error(current, "Parameter name is not found!")
			}
			paramName = false
		} else if paramType {
			if current.Type != fract.TypeDataType {
				fract.Error(current, "Parameter datatype is not found!")
			}
			paramType = false
		} else {
			if current.Type != fract.TypeComma {
				fract.Error(current, "Comma is not found!")
			}
			function.Parameters.Vals = append(function.Parameters.Vals, objects.Parameter{
				Name: paramList.Vals[index-2].(objects.Token).Value,
				Type: paramList.Vals[index-1].(objects.Token).Value,
			})
			paramName = true
			paramType = true
		}
	}

	if !paramName && paramType {
		fract.Error(paramList.Vals[len(paramList.Vals)-1].(objects.Token),
			"Parameter datatype is not found!")
	}

	i.funcs.Vals = append(i.funcs.Vals, function)

	i.skipBlock()
}
