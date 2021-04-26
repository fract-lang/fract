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

		modulus := -1
		multiplyOrDivive := -1
		binaryOrAnd := -1
		additionOrSubtraction := -1

		if token.Value == grammar.TokenPercent { // Modulus.
			if modulus == -1 {
				modulus = index
			}
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

		if modulus != -1 {
			if modulus == len(tokens)-1 {
				fract.Error(tokens[modulus], "Operator defined, but for what?")
			}
			return modulus
		} else if multiplyOrDivive != -1 {
			if multiplyOrDivive == len(tokens)-1 {
				fract.Error(tokens[multiplyOrDivive], "Operator defined, but for what?")
			}
			return multiplyOrDivive
		} else if binaryOrAnd != -1 {
			if binaryOrAnd == len(tokens)-1 {
				fract.Error(tokens[binaryOrAnd], "Operator defined, but for what?")
			}
			return binaryOrAnd
		} else if additionOrSubtraction != -1 {
			if additionOrSubtraction == len(tokens)-1 {
				fract.Error(tokens[additionOrSubtraction], "Operator defined, but for what?")
			}
			return additionOrSubtraction
		}
	}

	return -1
}
