/*
	processFunction Function.
*/

package interpreter

import (
	"../fract"
	"../fract/name"
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

	//paramList := tokens.Sublist(2, index-2)
	i.index++
	i.funcs.Vals = append(i.funcs.Vals, objects.Function{
		Name:  _name.Value,
		Start: i.index,
	})
	i.skipBlock()
}
