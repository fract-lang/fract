/*
	processVariableDefinition Function
*/

package interpreter

import (
	"../fract"
	"../fract/dt"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processVariable Process variable defination.
// tokens Tokens to process.
func (i *Interpreter) processVariableDefinition(tokens *vector.Vector) {
	var variable objects.Variable

	// Name is not defined?
	if len(tokens.Vals) < 2 {
		first := tokens.Vals[0].(objects.Token)
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Name is not found!")
	}

	_name := tokens.Vals[1].(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	// Name is already defined?
	if i.checkName(_name.Value) {
		fract.Error(_name, "Already defined this name!: "+_name.Value)
	}
	// Data type is not defined?
	if len(tokens.Vals) < 3 {
		fract.ErrorCustom(_name.File.Path, _name.Line, _name.Column+len(_name.Value),
			"Data type is not found!")
	}

	dataType := tokens.Vals[2].(objects.Token)
	// Data type is not data type token?
	if dataType.Type != fract.TypeDataType {
		fract.Error(dataType, "This is not a data type!")
	}
	// Setter is not defined?
	if len(tokens.Vals) < 4 {
		fract.ErrorCustom(dataType.File.Path, dataType.Line, dataType.Column+len(dataType.Value),
			"Setter is not found!")
	}

	setter := tokens.Vals[3].(objects.Token)
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.Setter {
		fract.Error(setter, "This is not a setter operator!"+setter.Value)
	}

	// Value is not defined?
	if len(tokens.Vals) < 5 {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	variable.Name = _name.Value
	variable.Type = dataType.Value
	valtokens := tokens.Sublist(4, len(tokens.Vals)-4)
	value := i.processValue(&valtokens)

	// Check value and data type compatibility.
	if dt.IsIntegerType(variable.Type) && value.Type != fract.VTInteger &&
		value.Type != fract.VTIntegerArray {
		fract.Error(setter, "Value and data type is not compatible!")
	}

	variable.Array = dt.TypeIsArray(value.Type)

	result, err := parser.ValueToTypeValue(variable.Array, variable.Type, value.Content)
	if err != "" {
		fract.ErrorCustom(setter.File.Path, setter.Line,
			setter.Column+len(setter.Value), err)
	}

	variable.Value = result

	// Set const state.
	variable.Const = tokens.Vals[0].(objects.Token).Value == grammar.KwConstVariable

	i.vars.Vals = append(i.vars.Vals, variable)
}
