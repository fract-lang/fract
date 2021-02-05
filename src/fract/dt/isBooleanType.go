/*
	IsBooleanType Function.
*/

package dt

import "../../grammar"

// IsBooleanType This type is a boolean.
// _type Type to check.
func IsBooleanType(_type string) bool {
	return _type == grammar.DtBoolean
}
