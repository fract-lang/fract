/*
	FindConditionOperator Function.
*/

package parser

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/grammar"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// FindConditionOperator Find next condition operator.
// tokens Tokens to search.
func FindConditionOperator(tokens vector.Vector) (int, string) {
	for index, current := range tokens.Vals {
		current := current.(objects.Token)
		if current.Type == fract.TypeOperator && (current.Value == grammar.Equals ||
			current.Value == grammar.NotEquals || current.Value == grammar.TokenGreat ||
			current.Value == grammar.TokenLess || current.Value == grammar.GreaterEquals ||
			current.Value == grammar.LessEquals) {
			return index, current.Value
		}
	}

	// Not found.
	return -1, ""
}
