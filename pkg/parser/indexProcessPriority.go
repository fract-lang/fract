package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// IndexProcessPriority Find index of priority operator.
// Returns index of operator if found, returns -1 if not.
//
// tokens Tokens to search.
func IndexProcessPriority(tokens vector.Vector) int {
	bracket := 0
	modulus := fract.TypeNone
	multiplyOrDivive := fract.TypeNone
	additionOrSubtraction := fract.TypeNone

	for index, _token := range tokens.Vals {
		_token := _token.(obj.Token)

		if _token.Type == fract.TypeBrace {
			if _token.Value == grammar.TokenLBracket ||
				_token.Value == grammar.TokenLBrace {
				bracket++
			} else if _token.Value == grammar.TokenRBracket ||
				_token.Value == grammar.TokenRBrace {
				bracket--
			}
		}

		if bracket > 0 {
			continue
		}

		// Exponentiation.
		if _token.Value == grammar.TokenCaret {
			return index
		} else if _token.Value == grammar.TokenPercent { // Modulus.
			if modulus == fract.TypeNone {
				modulus = index
			}
		} else if _token.Value == grammar.TokenStar ||
			_token.Value == grammar.TokenSlash ||
			_token.Value == grammar.TokenBackslash ||
			_token.Value == grammar.IntegerDivision ||
			_token.Value == grammar.IntegerDivideWithBigger { // Multiply or division.
			if multiplyOrDivive == fract.TypeNone {
				multiplyOrDivive = index
			}
		} else if _token.Value == grammar.TokenPlus ||
			_token.Value == grammar.TokenMinus { // Addition or subtraction.
			if additionOrSubtraction == fract.TypeNone {
				additionOrSubtraction = index
			}
		}
	}

	if modulus != fract.TypeNone {
		if modulus == len(tokens.Vals)-1 {
			fract.Error(tokens.Vals[modulus].(obj.Token),
				"Operator defined, but for what?")
		}
		return modulus
	} else if multiplyOrDivive != fract.TypeNone {
		if multiplyOrDivive == len(tokens.Vals)-1 {
			fract.Error(tokens.Vals[multiplyOrDivive].(obj.Token),
				"Operator defined, but for what?")
		}
		return multiplyOrDivive
	} else if additionOrSubtraction != fract.TypeNone {
		if additionOrSubtraction == len(tokens.Vals)-1 {
			fract.Error(tokens.Vals[additionOrSubtraction].(obj.Token),
				"Operator defined, but for what?")
		}
		return additionOrSubtraction
	}

	return fract.TypeNone
}
