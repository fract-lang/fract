package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

// IndexProcessPriority find index of priority operator and
// returns index of operator if found, returns -1 if not.
func IndexProcessPriority(tokens []objects.Token) int {
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

		switch token.Value {
		case grammar.TokenPercent: // Modulus.
			return index
		case grammar.TokenStar,
			grammar.TokenSlash,
			grammar.TokenBackslash,
			grammar.IntegerDivision,
			grammar.IntegerDivideWithBigger: // Multiply or division.
			if multiplyOrDivive == -1 {
				multiplyOrDivive = index
			}
		case grammar.TokenPlus,
			grammar.TokenMinus: // Addition or subtraction.
			if additionOrSubtraction == -1 {
				additionOrSubtraction = index
			}
		case grammar.TokenAmper,
			grammar.TokenVerticalBar:
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
