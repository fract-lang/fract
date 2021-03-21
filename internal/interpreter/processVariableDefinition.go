/*
	processVariableDefinition Function
*/

package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processVariable Process variable defination.
// tokens Tokens to process.
// protected Protected?
func (i *Interpreter) processVariableDefinition(tokens vector.Vector,
	protected bool) {
	tokenLen := len(tokens.Vals)

	// Name is not defined?
	if tokenLen < 2 {
		first := tokens.Vals[0].(obj.Token)
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value),
			"Name is not found!")
	}

	_name := tokens.Vals[1].(obj.Token)

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	} else if index := i.varIndexByName(_name.Value); index != -1 { // Name is already defined?
		fract.Error(_name, "Variable already defined in this name at line: "+
			fmt.Sprint(i.vars[index].Line))
	}

	// Data type is not defined?
	if tokenLen < 3 {
		fract.ErrorCustom(_name.File, _name.Line, _name.Column+len(_name.Value),
			"Setter is not found!")
	}

	setter := tokens.Vals[2].(obj.Token)
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.TokenEquals &&
		setter.Value != grammar.Input {
		fract.Error(setter, "This is not a setter operator!: "+setter.Value)
	}

	// Value is not defined?
	if tokenLen < 4 {
		fract.ErrorCustom(setter.File, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	var value obj.Value
	if setter.Value == grammar.TokenEquals { // =
		value = i.processValue(tokens.Sublist(3, tokenLen-3))
		if value.Content == nil {
			fract.Error(tokens.Vals[3].(obj.Token), "Invalid value!")
		}
	} else { // <<
		value = i.processInput(*tokens.Sublist(3, tokenLen-3))
	}

	i.vars = append(i.vars, obj.Variable{
		Name:      _name.Value,
		Value:     value,
		Line:      _name.Line,
		Const:     tokens.Vals[0].(obj.Token).Value == grammar.KwConstVariable,
		Protected: protected,
	})
}
