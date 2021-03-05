/*
	processVariableSet Function.
*/

package interpreter

import (
	"../fract"
	"../fract/arithmetic"
	"../fract/name"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// ProcessVariableSet Process variable set statement.
// tokens Tokens to process.
func (i *Interpreter) processVariableSet(tokens *vector.Vector) {
	_name := tokens.Vals[0].(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	index := name.VarIndexByName(i.vars, _name.Value)
	if index == -1 {
		fract.Error(_name, "Name is not defined!: "+_name.Value)
	}

	var setIndex int64 = -1
	variable := i.vars.Vals[index].(objects.Variable)
	setter := tokens.Vals[1].(objects.Token)

	// Array setter?
	if setter.Type == fract.TypeBrace && setter.Value == grammar.TokenLBracket {
		// Variable is not array?
		if !variable.Value.Array {
			fract.Error(setter, "Variable is not array!")
		}
		// Find close bracket.
		for cindex := 2; cindex < len(tokens.Vals); cindex++ {
			current := tokens.Vals[cindex].(objects.Token)
			if current.Type == fract.TypeBrace && current.Value == grammar.TokenRBracket {
				valueList := tokens.Sublist(2, cindex-2)
				// Index value is empty?
				if len(valueList.Vals) == 0 {
					fract.Error(setter, "Index is not defined!")
				}
				position, err := arithmetic.ToInt64(i.processValue(valueList).Content[0])
				if err != nil {
					fract.Error(setter, "Value out of range!")
				}
				if position < 0 || position >= int64(len(variable.Value.Content)) {
					fract.Error(setter, "Index is out of range!")
				}
				setIndex = position
				tokens.RemoveRange(1, cindex)
				setter = tokens.Vals[1].(objects.Token)
				break
			}
		}
	}

	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.Setter {
		fract.Error(setter, "This is not a setter operator!"+setter.Value)
	}

	// Value are not defined?
	if len(tokens.Vals) < 3 {
		fract.ErrorCustom(setter.File.Path, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	value := i.processValue(tokens.Sublist(2, len(tokens.Vals)-2))

	// Check const state
	if variable.Const {
		fract.Error(setter, "Values is can not changed of const defines!")
	}

	if setIndex != -1 {
		if !variable.Value.Array {
			fract.Error(_name, "This variable is not array!")
		}
		if value.Array {
			fract.Error(setter, "Array is cannot set as indexed value!")
		}
		variable.Value.Content[setIndex] = value.Content[0]
	} else {
		variable.Value = value
	}

	i.vars.Vals[index] = variable
}
