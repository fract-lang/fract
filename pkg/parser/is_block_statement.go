package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// IsBlockStatement returns true if tokens is block start, return false if not.
func IsBlockStatement(tokens []objects.Token) bool {
	if tokens[0].Type == fract.TypeMacro { // Remove macro token.
		tokens = tokens[1:]
	}

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
