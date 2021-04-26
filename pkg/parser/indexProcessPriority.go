package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// IndexProcessPriority Find index of priority operator.
// Returns index of operator if found, returns -1 if not.
//
// tokens Tokens to search.
func IndexProcessPriority(tokens []obj.Token) int {
	bracket := 0

	for index, token := range tokens {
		if token.Type == fract.TypeBrace {
			if token.Value == grammar.TokenLBracket ||
				token.Value == grammar.TokenLBrace ||
				token.Value == grammar.TokenLParenthes {
				bracket++
			} else {
				bracket--
			}
		}

		if bracket > 0 {
			continue
		}

		// Exponentiation or shifts.
		if token.Value == grammar.LeftBinaryShift || token.Value == grammar.RightBinaryShift ||
			token.Value == grammar.Exponentiation {
			return index
		}

		multiplyOrDivive := -1
		binaryOrAnd := -1
		additionOrSubtraction := -1

		if token.Value == grammar.TokenPercent { // Modulus.
			return index
		} else if token.Value == grammar.TokenStar ||
			token.Value == grammar.TokenSlash ||
			token.Value == grammar.TokenBackslash ||
			token.Value == grammar.IntegerDivision ||
			token.Value == grammar.IntegerDivideWithBigger { // Multiply or division.
			if multiplyOrDivive == -1 {
				multiplyOrDivive = index
			}
		} else if token.Value == grammar.TokenPlus ||
			token.Value == grammar.TokenMinus { // Addition or subtraction.
			if additionOrSubtraction == -1 {
				additionOrSubtraction = index
			}
		} else if token.Value == grammar.TokenAmper ||
			token.Value == grammar.TokenVerticalBar {
			if binaryOrAnd == -1 {
				binaryOrAnd = index
			}
		}

		if multiplyOrDivive != -1 {
			return multiplyOrDivive
		} else if binaryOrAnd != -1 {
			return binaryOrAnd
		} else if additionOrSubtraction != -1 {
			return additionOrSubtraction
		}
	}

	return -1
}
