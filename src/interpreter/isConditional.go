/*
	isConditional Function.
*/

package interpreter

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// isConditional Expression is conditional?
// tokens Tokens to check?
func (i *Interpreter) isConditional(tokens *vector.Vector) bool {
	// Search conditional expression.
	for index := 0; index < tokens.Len(); index++ {
		current := tokens.At(index).(objects.Token)
		if current.Type == fract.TypeOperator &&
			(current.Value == grammar.TokenAmper || current.Value == grammar.TokenVerticalBar ||
				current.Value == grammar.TokenEquals || current.Value == grammar.NotEquals ||
				current.Value == grammar.TokenGreat || current.Value == grammar.TokenLess) {
			return true
		}
	}

	return false
}
