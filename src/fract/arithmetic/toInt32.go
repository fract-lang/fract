/*
	ToInt32 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToInt32 String to 32bit integer.
// value Value to parse.
func ToInt32(value string) (int32, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(result), err
}
