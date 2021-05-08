package interpreter

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// appendVariable to source from tokens.
func (i *Interpreter) appendVariable(constant, protected bool, tokens []objects.Token) {
	_name := tokens[0]

	if strings.Contains(_name.Value, grammar.TokenDot) {
		fract.Error(_name, "Names is cannot include dot!")
	} else if _name.Value == grammar.TokenUnderscore {
		fract.Error(_name, "Ignore operator is cannot be variable name!")
	}

	// Name is already defined?
	if index, _ := i.varIndexByName(_name); index != -1 {
		fract.Error(_name, "Variable already defined in this name at line: "+fmt.Sprint(i.variables[index].Line))
	}

	tokensLen := len(tokens)

	// Setter is not defined?
	if tokensLen < 2 {
		fract.ErrorCustom(_name.File, _name.Line, _name.Column+len(_name.Value), "Setter is not found!")
	}

	setter := tokens[1]
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != grammar.TokenEquals {
		fract.Error(setter, "This is not a setter operator!: "+setter.Value)
	}

	// Value is not defined?
	if tokensLen < 3 {
		fract.ErrorCustom(setter.File, setter.Line, setter.Column+len(setter.Value), "Value is not defined!")
	}

	value := i.processValue(*vector.Sublist(tokens, 2, tokensLen-2))
	if value.Content == nil {
		fract.Error(tokens[2], "Invalid value!")
	}

	i.variables = append(i.variables,
		objects.Variable{
			Name:      _name.Value,
			Value:     value,
			Line:      _name.Line,
			Const:     constant,
			Protected: protected,
		})
}

func (i *Interpreter) processVariableDefinition(tokens []objects.Token, protected bool) {
	// Name is not defined?
	if len(tokens) < 2 {
		first := tokens[0]
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value), "Name is not found!")
	}

	constant := tokens[0].Value == grammar.KwConstVariable
	pre := tokens[1]

	if pre.Type == fract.TypeName {
		i.appendVariable(constant, protected, tokens[1:])
	} else if pre.Type == fract.TypeBrace && pre.Value == grammar.TokenLParenthes {
		tokens = tokens[2 : len(tokens)-1]
		last := 0
		bracket := 0
		for index, token := range tokens {
			if token.Type == fract.TypeBrace {
				if token.Value == grammar.TokenLBrace ||
					token.Value == grammar.TokenLBracket ||
					token.Value == grammar.TokenLParenthes {
					bracket++
				} else {
					bracket--
				}
			}

			if bracket > 0 {
				continue
			}

			if token.Type == fract.TypeComma {
				i.appendVariable(constant, protected, tokens[last:index])
				last = index + 1
			}
		}

		if len(tokens) != last {
			i.appendVariable(constant, protected, tokens[last:])
		}
	} else {
		fract.Error(pre, "Invalid syntax!")
	}
}
