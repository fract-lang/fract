package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

func (i *Interpreter) processDelete(tokens []objects.Token) {
	tokenLen := len(tokens)
	// Value is not defined?
	if tokenLen < 2 {
		first := tokens[0]
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value), "Value is not defined!")
	}

	comma := false
	for index := 1; index < tokenLen; index++ {
		current := tokens[index]
		if comma {
			if current.Type != fract.TypeComma {
				fract.Error(current, "Comma is not found!")
			}
			comma = false
			continue
		}
		// Token is not a deletable object?
		if current.Type != fract.TypeName {
			fract.Error(current, "This is not deletable object!")
		}
		if index < tokenLen-1 {
			next := tokens[index+1]
			if next.Type == fract.TypeBrace && next.Value == "(" {
				nnext := tokens[index+2]
				if !(nnext.Type == fract.TypeBrace && nnext.Value == ")") {
					fract.Error(nnext, "Invalid syntax!")
				}
				index += 2
				position, source := i.functionIndexByName(current)
				// Name is not defined?
				if position == -1 {
					fract.Error(current, "Name is not defined!")
				}
				// Protected?
				if source.functions[position].Protected {
					fract.Error(current, "Protected objects cannot be deleted manually from memory!")
				}
				source.functions = append(source.functions[:position], source.functions[position+1:]...)
				comma = true
				continue
			}
		}
		position, source := i.varIndexByName(current)
		// Name is not defined?
		if position == -1 {
			fract.Error(current, "Name is not defined!")
		}
		// Protected?
		if source.variables[position].Protected {
			fract.Error(current, "Protected objects cannot be deleted manually from memory!")
		}
		source.variables = append(source.variables[:position], source.variables[position+1:]...)
		comma = true
	}
}
