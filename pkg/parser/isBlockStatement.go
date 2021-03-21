package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// IsBlockStatement Statement is block?
// tokens Tokens of statement.
func IsBlockStatement(tokens vector.Vector) bool {
	first := tokens.Vals[0].(objects.Token)
	if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
		first.Type == fract.TypeFunction {
		return true
	} else if first.Type == fract.TypeProtected {
		if len(tokens.Vals) > 1 {
			if second := tokens.Vals[1].(objects.Token); second.Type == fract.TypeFunction {
				return true
			}
		}
	}
	return false
}
