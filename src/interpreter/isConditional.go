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
func (i *Interpreter) isConditional(tokens vector.Vector) bool {
	// Search conditional expression.
	for index := range tokens.Vals {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeOperator &&
			(current.Value == grammar.TokenAmper || current.Value == grammar.TokenVerticalBar ||
				current.Value == grammar.TokenEquals || current.Value == grammar.NotEquals ||
				current.Value == grammar.TokenGreat || current.Value == grammar.TokenLess ||
				current.Value == grammar.GreaterEquals || current.Value == grammar.LessEquals) {
			return true
		}
	}

	return false
}
