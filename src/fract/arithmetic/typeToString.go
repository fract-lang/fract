/*
	TypeToString Function.
*/

package arithmetic

import (
	"strings"

	fract ".."
	"../../grammar"
)

// TypeToString Parse type to string.
// _type Type.
// value Value to parse.
func TypeToString(_type int, value interface{}) string {
	if _type == fract.VTFloat {
		return FloatToString(value)
	}
	str := FloatToString(value)
	return str[:strings.Index(str, grammar.TokenDot)]
}
