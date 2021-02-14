/*
	ToFloat32 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToFloat32 String to float.
// value Value to parse.
func ToFloat32(value string) (float32, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(result), err
}
