/*
	GENERIC
*/

package arithmetic

import (
	"strings"

	fract ".."
	"../../grammar"
)

// IsNegative Is negative number?
// value Value to check.
func IsNegative(value string) bool {
	return strings.HasPrefix(value, grammar.TokenMinus)
}

// IsNumeric Char is numeric?
// char Char to check.
func IsNumeric(char byte) bool {
	return char == '0' ||
		char == '1' ||
		char == '2' ||
		char == '3' ||
		char == '4' ||
		char == '5' ||
		char == '6' ||
		char == '7' ||
		char == '8' ||
		char == '9'
}

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
