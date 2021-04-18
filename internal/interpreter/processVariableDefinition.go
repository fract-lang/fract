/*
	processVariableDefinition Function
*/

package interpreter

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processVariable Process variable defination.
// tokens Tokens to process.
// protected Protected?
func (i *Interpreter) processVariableDefinition(tokens []obj.Token, protected bool) {
	// Name is not defined?
	if len(tokens) < 2 {
		first := tokens[0]
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value),
			"Name is not found!")
	}

	_const := tokens[0].Value == grammar.KwConstVariable

	appendVariable := func(avtokens []obj.Token) {
		_name := avtokens[0]
		if strings.Contains(_name.Value, grammar.TokenDot) {
			fract.Error(_name, "Names is cannot include dot!")
		}

		if _name.Value == grammar.TokenUnderscore {
			fract.Error(_name, "Ignore operator is cannot be variable name!")
		}

		// Name is already defined?
		if index, _ := i.varIndexByName(_name); index != -1 {
			fract.Error(_name, "Variable already defined in this name at line: "+
				fmt.Sprint(i.variables[index].Line))
		}

		tokensLen := len(avtokens)

		// Setter is not defined?
		if tokensLen < 2 {
			fract.ErrorCustom(_name.File, _name.Line, _name.Column+len(_name.Value),
				"Setter is not found!")
		}

		setter := avtokens[1]
		/*// Setter is not a setter operator?
		if setter.Type != fract.TypeOperator && setter.Value != grammar.TokenEquals {
			fract.Error(setter, "This is not a setter operator!: "+setter.Value)
		}*/

		// Value is not defined?
		if tokensLen < 3 {
			fract.ErrorCustom(setter.File, setter.Line, setter.Column+len(setter.Value),
				"Value is not defined!")
		}

		value := i.processValue(vector.Sublist(avtokens, 2, tokensLen-2))
		if value.Content == nil {
			fract.Error(avtokens[2], "Invalid value!")
		}

		i.variables = append(i.variables, obj.Variable{
			Name:      _name.Value,
			Value:     value,
			Line:      _name.Line,
			Const:     _const,
			Protected: protected,
		})
	}

	pre := tokens[1]

	if pre.Type == fract.TypeName {
		appendVariable(tokens[1:])
	} else if pre.Type == fract.TypeBrace && pre.Value == grammar.TokenLParenthes {
		tokens = tokens[2 : len(tokens)-1]
		last := 0
		for index, token := range tokens {
			if token.Type == fract.TypeComma {
				appendVariable(tokens[last:index])
				last = index + 1
			}
		}
		if len(tokens) != last {
			appendVariable(tokens[last:])
		}
	} else {
		fract.Error(pre, "Invalid syntax!")
	}
}
