package interpreter

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// findConditionOperator Find next condition operator.
// tokens Tokens to search.
func findConditionOperator(tokens *vector.Vector) (int, string) {
	for index := 0; index < tokens.Len(); index++ {
		current := tokens.At(index).(objects.Token)
		if current.Type == fract.TypeOperator && (current.Value == grammar.TokenEquals ||
			current.Value == grammar.NotEquals) {
			return index, current.Value
		}
	}

	// Not found.
	return -1, ""
}
