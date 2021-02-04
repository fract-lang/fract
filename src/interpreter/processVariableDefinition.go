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
	"../utilities/vector"
)

// processVariable Process variable defination.
// tokens Tokens to process.
func (i *Interpreter) processVariableDefinition(tokens *vector.Vector) {
	var variable objects.Variable

	/* Check size. */
	if tokens.Len() < 5 {
		fract.Error(tokens.Last().(objects.Token), "")
	}

	_name := tokens.At(1).(objects.Token)

	if name.VarIndexByName(i.vars, _name.Value) != -1 {
		fract.Error(_name, "Already exist variable in this name!: "+_name.Value)
	}

	dataType := tokens.At(2).(objects.Token)

	// Data type is not data type token?
	if dataType.Type != fract.TypeDataType {
		fract.Error(dataType, "This is not a data type!")
	}

	setter := tokens.At(3).(objects.Token)
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.Setter {
		fract.Error(setter, "This is not a setter operator("+grammar.Setter+")!")
	}

	variable.Name = _name.Value
	variable.Type = dataType.Value
	valtokens := tokens.Sublist(4, tokens.Len()-4)
	value := i.processValue(&valtokens)

	// Check value and data type compatibility.
	if dt.IsIntegerType(variable.Type) && value.Type != fract.VTInteger {
		fract.Error(setter, "Value and data type is not compatible!")
	}

	variable.Value = value.Content

	i.vars.Append(variable)
}
