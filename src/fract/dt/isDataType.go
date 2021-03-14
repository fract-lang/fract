/*
	IsDataType Function.
*/

package dt

import (
	"github.com/fract-lang/src/grammar"
)

// IsDataType Check value is a data type or not?
// value Value to check.
func IsDataType(value string) bool {
	return value == grammar.DtInt8 ||
		value == grammar.DtInt16 ||
		value == grammar.DtInt32 ||
		value == grammar.DtInt64 ||
		value == grammar.DtUInt8 ||
		value == grammar.DtUInt16 ||
		value == grammar.DtUInt32 ||
		value == grammar.DtUInt64 ||
		value == grammar.DtBoolean
}
