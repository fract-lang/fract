/*
	processVariableDefinition Function
*/

package interpreter

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utils/vector"
)

// processVariable Process variable defination.
// tokens Tokens to process.
func (i *Interpreter) processVariableDefinition(tokens vector.Vector) {
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
	} else if i.varIndexByName(_name.Value) != -1 { // Name is already defined?
		fract.Error(_name, "Already defined this name!: "+_name.Value)
	}

	// Data type is not defined?
	if tokenLen < 3 {
		fract.ErrorCustom(_name.File.Path, _name.Line, _name.Column+len(_name.Value),
			"Setter is not found!")
	}

	setter := tokens.Vals[2].(objects.Token)
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.TokenEquals &&
		setter.Value != grammar.Input {
		fract.Error(setter, "This is not a setter operator!: "+setter.Value)
	}

	// Value is not defined?
	if tokenLen < 4 {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	variable.Name = _name.Value
	var value objects.Value
	if setter.Value == grammar.TokenEquals { // =
		value = i.processValue(tokens.Sublist(3, tokenLen-3))
	} else { // <<
		value = i.processInput(*tokens.Sublist(3, tokenLen-3))
	}

	variable.Value = value

	// Set const state.
	variable.Const = tokens.Vals[0].(objects.Token).Value == grammar.KwConstVariable

	i.vars.Vals = append(i.vars.Vals, variable)
}
