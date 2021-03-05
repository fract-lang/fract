/*
	processVariableDefinition Function
*/

package interpreter

import (
	"../fract"
	"../fract/name"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// processVariable Process variable defination.
// tokens Tokens to process.
func (i *Interpreter) processVariableDefinition(tokens *vector.Vector) {
	var variable objects.Variable

	tokenLen := len(tokens.Vals)

	// Name is not defined?
	if tokenLen < 2 {
		first := tokens.Vals[0].(objects.Token)
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Name is not found!")
	}

	_name := tokens.Vals[1].(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	} else if name.VarIndexByName(i.vars, _name.Value) != -1 { // Name is already defined?
		fract.Error(_name, "Already defined this name!: "+_name.Value)
	}

	// Data type is not defined?
	if tokenLen < 3 {
		fract.ErrorCustom(_name.File.Path, _name.Line, _name.Column+len(_name.Value),
			"Setter is not found!")
	}

	setter := tokens.Vals[2].(objects.Token)
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.Setter {
		fract.Error(setter, "This is not a setter operator!: "+setter.Value)
	}

	// Value is not defined?
	if tokenLen < 4 {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	variable.Name = _name.Value
	value := i.processValue(tokens.Sublist(3, tokenLen-3))

	variable.Value = value

	// Set const state.
	variable.Const = tokens.Vals[0].(objects.Token).Value == grammar.KwConstVariable

	i.vars.Vals = append(i.vars.Vals, variable)
}
