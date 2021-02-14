/*
	CheckFloat Function.
*/

package arithmetic

import (
	"strings"

	"../../grammar"
)

// CheckFloat Check float value validate.
// value Value to check.
func CheckFloat(value string) bool {
	return len(value[strings.Index(value, grammar.TokenDot)+1:]) <= 6
}
