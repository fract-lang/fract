/*
	ToUInt16 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToUInt16 String to 16bit unsigned integer.
// value Value to parse.
func ToUInt16(value string) (uint16, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(result), err
}
