/*
	processVariableDefinition Function
*/

package interpreter

import (
	"../fract"
	"../fract/dt"
	"../fract/name"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processVariable Process variable defination.
// tokens Tokens to process.
func (i *Interpreter) processVariableDefinition(tokens *vector.Vector) {
	var variable objects.Variable
	length := tokens.Len()

	// Name is not defined?
	if length < 2 {
		first := tokens.First().(objects.Token)
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Name is not found!")
	}

	_name := tokens.At(1).(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	// Name is already defined?
	if name.VarIndexByName(i.vars, _name.Value) != -1 {
		fract.Error(_name, "Already defined this name!: "+_name.Value)
	}

	// Data type is not defined?
	if length < 3 {
		fract.ErrorCustom(_name.File.Path, _name.Line, _name.Column+len(_name.Value),
			"Data type is not found!")
	}

	dataType := tokens.At(2).(objects.Token)

	// Data type is not data type token?
	if dataType.Type != fract.TypeDataType {
		fract.Error(dataType, "This is not a data type!")
	}

	// Setter is not defined?
	if length < 4 {
		fract.ErrorCustom(dataType.File.Path, dataType.Line, dataType.Column+len(dataType.Value),
			"Setter is not found!")
	}

	setter := tokens.At(3).(objects.Token)
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.Setter {
		fract.Error(setter, "This is not a setter operator!"+setter.Value)
	}

	// Value is not defined?
	if length < 5 {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	variable.Name = _name.Value
	variable.Type = dataType.Value
	valtokens := tokens.Sublist(4, tokens.Len()-4)
	value := i.processValue(&valtokens)

	// Check value and data type compatibility.
	if dt.IsIntegerType(variable.Type) && value.Type != fract.VTInteger {
		fract.Error(setter, "Value and data type is not compatible!")
	}

	result, err := parser.ValueToTypeValue(variable.Type, value.Content)
	if err != "" {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value), err)
	}
	variable.Value = result

	i.vars.Append(variable)
}
