package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

// FindConditionOperator return next condition operator.
func FindConditionOperator(tokens []objects.Token) (int, string) {
	for index, current := range tokens {
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
