/*
	ToUInt32 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToUInt32 String to 32bit integer.
// value Value to parse.
func ToUInt32(value string) (uint32, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(result), err
}
