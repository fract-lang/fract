/*
	processVariableSet Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/fract/arithmetic"
	"github.com/fract-lang/fract/src/grammar"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// ProcessVariableSet Process variable set statement.
// tokens Tokens to process.
func (i *Interpreter) processVariableSet(tokens vector.Vector) {
	_name := tokens.Vals[0].(objects.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	index := i.varIndexByName(_name.Value)
	if index == -1 {
		fract.Error(_name, "Variable is not defined in this name!: "+_name.Value)
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
				if valueList.Vals == nil {
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
	if setter.Type != fract.TypeOperator && setter.Value != grammar.TokenEquals &&
		setter.Value != grammar.Input {
		fract.Error(setter, "This is not a setter operator!"+setter.Value)
	}

	// Value are not defined?
	if len(tokens.Vals) < 3 {
		fract.ErrorCustom(setter.File, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	var value objects.Value
	if setter.Value == grammar.TokenEquals { // =
		value = i.processValue(tokens.Sublist(2, len(tokens.Vals)-2))
	} else { // <<
		value = i.processInput(*tokens.Sublist(2, len(tokens.Vals)-2))
	}

	// Check const state
	if variable.Const {
		fract.Error(setter, "Values is can not changed of const defines!")
	}

	if setIndex != -1 {
		if value.Array {
			fract.Error(setter, "Array is cannot set as indexed value!")
		}
		variable.Value.Content[setIndex] = value.Content[0]
	} else {
		variable.Value = value
	}

	i.vars.Vals[index] = variable
}
