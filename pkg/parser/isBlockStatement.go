package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// IsBlockStatement Statement is block?
// tokens Tokens of statement.
func IsBlockStatement(tokens []obj.Token) bool {
	first := tokens[0]
	if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
		first.Type == fract.TypeFunction || first.Type == fract.TypeTry {
		return true
	} else if first.Type == fract.TypeProtected {
		if len(tokens) > 1 {
			if tokens[1].Type == fract.TypeFunction {
				return true
			}
		}
	}
	return false
}
