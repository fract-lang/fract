/*
	IsIntegerType Function.
*/

package dt

import "../../grammar"

// IsIntegerType This type is a integer.
// _type Type to check.
func IsIntegerType(_type string) bool {
	return _type == grammar.DtInt8 ||
		_type == grammar.DtInt32 ||
		_type == grammar.DtInt64 ||
		_type == grammar.DtUInt8 ||
		_type == grammar.DtUInt32 ||
		_type == grammar.DtUInt64
}
