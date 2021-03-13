/*
	processDelete Function.
*/

package interpreter

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utils/vector"
)

// processDelete Process delete keyword.
// tokens Tokens to process.
func (i *Interpreter) processDelete(tokens vector.Vector) {
	tokenLen := len(tokens.Vals)

	// Value is not defined?
	if tokenLen < 2 {
		first := tokens.Vals[0].(objects.Token)
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Value is not found!")
	}

	comma := false
	for index := 1; index < tokenLen; index++ {
		current := tokens.Vals[index].(objects.Token)

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
			next := tokens.Vals[index+1].(objects.Token)
			if next.Type == fract.TypeBrace && next.Value == grammar.TokenLParenthes {
				nnext := tokens.Vals[index+2].(objects.Token)
				if !(nnext.Type == fract.TypeBrace && nnext.Value == grammar.TokenRParenthes) {
					fract.Error(nnext, "Invalid syntax!")
				}
				index += 2

				position := i.functionIndexByName(current.Value)

				// Name is not defined?
				if position == -1 {
					fract.Error(current, "Name is not defined!")
				}

				i.funcs.RemoveRange(position, 1)
				comma = true
				continue
			}
		}

		position := i.varIndexByName(current.Value)

		// Name is not defined?
		if position == -1 {
			fract.Error(current, "Name is not defined!")
		}

		i.vars.RemoveRange(position, 1)
		comma = true
	}
}
