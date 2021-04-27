package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// IsBlockStatement Statement is block?
// tokens Tokens of statement.
func IsBlockStatement(tokens []obj.Token) bool {
	switch tokens[0].Type {
	case fract.TypeIf,
		fract.TypeLoop,
		fract.TypeFunction,
		fract.TypeTry:
		return true
	case fract.TypeProtected:
		if len(tokens) > 1 {
			if tokens[1].Type == fract.TypeFunction {
				return true
			}
		}
	}
	return false
}
