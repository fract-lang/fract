package parser

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// IndexProcessPriority Find index of priority operator.
// Returns index of operator if found, returns -1 if not.
//
// tokens Tokens to search.
func IndexProcessPriority(tokens vector.Vector) int {
	bracket := 0
	modulus := -1
	multiplyOrDivive := -1
	additionOrSubtraction := -1

	for index := range tokens.Vals {
		_token := tokens.Vals[index].(objects.Token)

		if _token.Type == fract.TypeBrace {
			if _token.Value == grammar.TokenLBracket || _token.Value == grammar.TokenLBrace {
				bracket++
			} else if _token.Value == grammar.TokenRBracket || _token.Value == grammar.TokenRBrace {
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
			if modulus == -1 {
				modulus = index
			}
		} else if _token.Value == grammar.TokenStar ||
			_token.Value == grammar.TokenSlash ||
			_token.Value == grammar.TokenBackslash ||
			_token.Value == grammar.IntegerDivision ||
			_token.Value == grammar.IntegerDivideWithBigger { // Multiply or division.
			if multiplyOrDivive == -1 {
				multiplyOrDivive = index
			}
		} else if _token.Value == grammar.TokenPlus ||
			_token.Value == grammar.TokenMinus { // Addition or subtraction.
			if additionOrSubtraction == -1 {
				additionOrSubtraction = index
			}
		}
	}

	if modulus != -1 {
		if modulus == len(tokens.Vals)-1 {
			fract.Error(tokens.Vals[modulus].(objects.Token),
				"Operator defined, but for what?")
		}
		return modulus
	} else if multiplyOrDivive != -1 {
		if multiplyOrDivive == len(tokens.Vals)-1 {
			fract.Error(tokens.Vals[multiplyOrDivive].(objects.Token),
				"Operator defined, but for what?")
		}
		return multiplyOrDivive
	} else if additionOrSubtraction != -1 {
		if additionOrSubtraction == len(tokens.Vals)-1 {
			fract.Error(tokens.Vals[additionOrSubtraction].(objects.Token),
				"Operator defined, but for what?")
		}
		return additionOrSubtraction
	}

	return -1
}
