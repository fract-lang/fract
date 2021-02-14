/*
	ToUInt64 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToUInt64 String to 64bit unsigned integer.
// value Value to parse.
func ToUInt64(value string) (uint64, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(result), err
}
