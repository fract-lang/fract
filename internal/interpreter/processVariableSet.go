/*
	processVariableSet Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// ProcessVariableSet Process variable set statement.
// tokens Tokens to process.
func (i *Interpreter) processVariableSet(tokens []obj.Token) {
	_name := tokens[0]

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	}

	index := i.varIndexByName(_name.Value)
	if index == -1 {
		fract.Error(_name, "Variable is not defined in this name!: "+_name.Value)
	}

	setIndex := -1
	variable := i.vars[index]
	setter := tokens[1]

	// Check const state
	if variable.Const {
		fract.Error(setter, "Values is cannot changed of constant defines!")
	}

	// Array setter?
	if setter.Type == fract.TypeBrace && setter.Value == grammar.TokenLBracket {
		// Variable is not array?
		if !variable.Value.Array {
			fract.Error(setter, "Variable is not array!")
		}
		// Find close bracket.
		for cindex := 2; cindex < len(tokens); cindex++ {
			current := tokens[cindex]
			if current.Type != fract.TypeBrace ||
				current.Value != grammar.TokenRBracket {
				continue
			}

			valueList := vector.Sublist(tokens, 2, cindex-2)
			// Index value is empty?
			if valueList == nil {
				fract.Error(setter, "Index is not defined!")
			}
			position, err := arithmetic.ToInt(i.processValue(valueList).Content[0])
			if err != nil {
				fract.Error(setter, "Value out of range!")
			}
			position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
			if position == -1 {
				fract.Error(setter, "Index is out of range!")
			}
			setIndex = position
			vector.RemoveRange(&tokens, 1, cindex)
			setter = tokens[1]
			break
		}
	}

	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.TokenEquals &&
		setter.Value != grammar.Input {
		fract.Error(setter, "This is not a setter operator!"+setter.Value)
	}

	// Value are not defined?
	if len(tokens) < 3 {
		fract.ErrorCustom(setter.File, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	var value obj.Value
	if setter.Value == grammar.TokenEquals { // =
		value = i.processValue(vector.Sublist(tokens, 2, len(tokens)-2))
		if value.Content == nil {
			fract.Error(tokens[2], "Invalid value!")
		}
	} else { // <<
		value = i.processInput(*vector.Sublist(tokens, 2, len(tokens)-2))
	}

	if setIndex != -1 {
		if value.Array {
			fract.Error(setter, "Array is cannot set as indexed value!")
		}
		variable.Value.Content[setIndex] = value.Content[0]
	} else {
		variable.Value = value
	}

	i.vars[index] = variable
}
