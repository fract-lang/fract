/*
	IsFloatType Function.
*/

package dt

import "../../grammar"

// IsFloatType This type is a float.
// _type Type to check.
func IsFloatType(_type string) bool {
	return _type == grammar.DtFloat32 ||
		_type == grammar.DtFloat64
}
