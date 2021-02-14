/*
	ToUInt8 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToUInt8 String to 8bit unsigned integer.
// value Value to parse.
func ToUInt8(value string) (uint8, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(result), err
}
