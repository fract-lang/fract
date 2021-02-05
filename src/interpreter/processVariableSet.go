/*
	processVariableSet Function.
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

// ProcessVariableSet Process variable set statement.
// tokens Tokens to process.
func (i *Interpreter) processVariableSet(tokens *vector.Vector) {
	_name := tokens.At(0).(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	index := name.VarIndexByName(i.vars, _name.Value)
	if index == -1 {
		fract.Error(_name, "Name is not defined!: "+_name.Value)
	}

	setter := tokens.At(1).(objects.Token)

	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.Setter {
		fract.Error(setter, "This is not a setter operator!"+setter.Value)
	}

	// Value are not defined?
	if tokens.Len() < 3 {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	variable := i.vars.At(index).(objects.Variable)
	valtokens := tokens.Sublist(2, tokens.Len()-2)
	value := i.processValue(&valtokens)

	// Check value and data type compatibility.
	if dt.IsIntegerType(variable.Type) && value.Type != fract.VTInteger {
		fract.Error(setter, "Value and data type is not compatible!")
	}

	result, err := parser.ValueToTypeValue(variable.Type, value.Content)
	if err != "" {
		fract.ErrorCustom(setter.File.Path, setter.Line,
			setter.Column+len(setter.Value), err)
	}
	variable.Value = result
	i.vars.Set(index, variable)
}
