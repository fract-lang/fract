/*
	processVariableSet Function.
*/

package interpreter

import (
	"../fract"
	"../fract/arithmetic"
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
	var setIndex int64 = -1
	variable := i.vars.At(index).(objects.Variable)

	setter := tokens.At(1).(objects.Token)

	// Array setter?
	if setter.Type == fract.TypeBrace && setter.Value == grammar.TokenLBracket {
		// Variable is not array?
		if !variable.Array {
			fract.Error(setter, "Variable is not array!")
		}
		// Find close bracket.
		for cindex := 2; cindex < tokens.Len(); cindex++ {
			current := tokens.At(cindex).(objects.Token)
			if current.Type == fract.TypeBrace && current.Value == grammar.TokenRBracket {
				valueList := tokens.Sublist(2, cindex-2)
				position, err := arithmetic.ToInt64(i.processValue(&valueList).Content[0])
				if err != nil {
					fract.Error(setter, "Value out of range!")
				}
				if position < 0 || position >= int64(len(variable.Value)) {
					fract.Error(setter, "Index is out of range!")
				}
				setIndex = position
				tokens.RemoveRange(1, cindex)
				setter = tokens.At(1).(objects.Token)
				break
			}
		}
	}

	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.Setter {
		fract.Error(setter, "This is not a setter operator!"+setter.Value)
	}

	// Value are not defined?
	if tokens.Len() < 3 {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	valtokens := tokens.Sublist(2, tokens.Len()-2)
	value := i.processValue(&valtokens)

	if variable.Array && !dt.TypeIsArray(value.Type) && setIndex == -1 {
		fract.Error(setter, "This variable is array, cannot set nonarray value!")
	} else if !variable.Array && dt.TypeIsArray(value.Type) {
		fract.Error(setter, "This variable is not array, cannot set array value!")
	}

	// Check value and data type compatibility.
	if dt.IsIntegerType(variable.Type) && value.Type != fract.VTInteger &&
		value.Type != fract.VTIntegerArray {
		fract.Error(setter, "Value and data type is not compatible!")
	}

	result, err := parser.ValueToTypeValue(variable.Array, variable.Type, value.Content)
	if err != "" {
		fract.ErrorCustom(setter.File.Path, setter.Line,
			setter.Column+len(setter.Value), err)
	}

	// Check const state
	if variable.Const {
		fract.Error(setter, "Values is can not changed of const defines!")
	}

	if setIndex != -1 {
		variable.Value[setIndex] = result[0]
	} else {
		variable.Value = result
	}
	i.vars.Set(index, variable)
}
