/*
	IsArray Function.
*/

package dt

import "../../utils/vector"

// IsArray Value is array?
// value Value to check.
func IsArray(value interface{}) bool {
	switch value.(type) {
	case vector.Vector:
		return true
	default:
		return false
	}
}
