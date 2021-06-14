package interpreter

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// Metadata of variable declaration.
type variableMetadata struct {
	constant  bool
	mutable   bool
	protected bool
}

// appendVariable to source from tokens.
func (i *Interpreter) appendVariable(md variableMetadata, tokens []objects.Token) {
	_name := tokens[0]
	if strings.Contains(_name.Value, ".") {
		fract.Error(_name, "Names is cannot include dot!")
	} else if _name.Value == "_" {
		fract.Error(_name, "Ignore operator is cannot be variable name!")
	}
	// Name is already defined?
	if line := i.definedName(_name); line != -1 {
		fract.Error(_name, "\""+_name.Value+"\" is already defined at line: "+fmt.Sprint(line))
	}
	tokensLen := len(tokens)
	// Setter is not defined?
	if tokensLen < 2 {
		fract.ErrorCustom(_name.File, _name.Line, _name.Column+len(_name.Value), "Setter is not found!")
	}
	setter := tokens[1]
	// Setter is not a setter operator?
	if setter.Type != fract.TypeOperator && setter.Value != "=" {
		fract.Error(setter, "This is not a setter operator: "+setter.Value)
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
			Const:     md.constant,
			Mutable:   md.mutable,
			Protected: md.protected,
		})
}

func (i *Interpreter) processVariableDefinition(tokens []objects.Token, protected bool) {
	// Name is not defined?
	if len(tokens) < 2 {
		first := tokens[0]
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value), "Name is not found!")
	}
	md := variableMetadata{
		constant:  tokens[0].Value == grammar.KwConstant,
		mutable:   tokens[0].Value == grammar.KwMut,
		protected: protected,
	}
	pre := tokens[1]
	if pre.Type == fract.TypeName {
		i.appendVariable(md, tokens[1:])
	} else if pre.Type == fract.TypeBrace && pre.Value == "(" {
		tokens = tokens[2 : len(tokens)-1]
		last := 0
		line := tokens[0].Line
		bracket := 0
		for index, token := range tokens {
			if token.Type == fract.TypeBrace {
				if token.Value == "{" || token.Value == "[" || token.Value == "(" {
					bracket++
				} else {
					bracket--
					line = token.Line
				}
			}
			if bracket > 0 {
				continue
			}
			if line < token.Line {
				i.appendVariable(md, tokens[last:index])
				last = index
				line = token.Line
			}
		}
		if len(tokens) != last {
			i.appendVariable(md, tokens[last:])
		}
	} else {
		fract.Error(pre, "Invalid syntax!")
	}
}
