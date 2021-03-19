/*
	IsFloatValue Function.
*/

package arithmetic

import (
	"strings"

	"github.com/fract-lang/fract/pkg/grammar"
)

// IsFloatValue Value is float?
// value Value to check.
func IsFloatValue(value string) bool {
	return strings.Index(value, grammar.TokenDot) != -1
}
