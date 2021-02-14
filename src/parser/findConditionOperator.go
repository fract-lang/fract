/*
	FindConditionOperator Function.
*/

package parser

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// FindConditionOperator Find next condition operator.
// tokens Tokens to search.
func FindConditionOperator(tokens *vector.Vector) (int, string) {
	for index := range tokens.Vals {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeOperator && (current.Value == grammar.TokenEquals ||
			current.Value == grammar.NotEquals || current.Value == grammar.TokenGreat ||
			current.Value == grammar.TokenLess || current.Value == grammar.GreaterEquals ||
			current.Value == grammar.LessEquals) {
			return index, current.Value
		}
	}

	// Not found.
	return -1, ""
}
