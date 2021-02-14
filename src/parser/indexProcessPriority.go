package parser

import (
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// IndexProcessPriority Find index of priority operator.
// Returns index of operator if found, returns -1 if not.
//
// tokens Tokens to search.
func IndexProcessPriority(tokens *vector.Vector) int {
	// Returns -1 if vector contains one value.
	if len(tokens.Vals) == 1 {
		return -1
	}

	/* Find exponentiation. */
	for index := range tokens.Vals {
		if tokens.Vals[index].(objects.Token).Value == grammar.TokenCaret {
			return index
		}
	}

	/* Find mod. */
	for index := range tokens.Vals {
		if tokens.Vals[index].(objects.Token).Value == grammar.TokenPercent {
			return index
		}
	}

	/* Find multipy or divide. */
	for index := range tokens.Vals {
		_token := tokens.Vals[index].(objects.Token)
		if _token.Value == grammar.TokenStar ||
			_token.Value == grammar.TokenSlash ||
			_token.Value == grammar.TokenReverseSlash ||
			_token.Value == grammar.IntegerDivision ||
			_token.Value == grammar.IntegerDivideWithBigger {
			return index
		}
	}

	/* Addition or subtraction. */
	for index := range tokens.Vals {
		_token := tokens.Vals[index].(objects.Token)
		if _token.Value == grammar.TokenPlus || _token.Value == grammar.TokenMinus {
			return index
		}
	}

	return -1
}
