package dt

import (
	fract ".."
)

// TypeIsArray Type is array?
// _type Type to check.
func TypeIsArray(_type int) bool {
	return _type == fract.VTIntegerArray || _type == fract.VTFloatArray
}
