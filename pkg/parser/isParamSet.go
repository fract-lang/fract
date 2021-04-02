/*
	IsParamSet Function.
*/

package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// IsParamSet Argument type is param set?
func IsParamSet(tokens []obj.Token) bool {
	return tokens[0].Type == fract.TypeName &&
		tokens[1].Value == grammar.TokenEquals
}
