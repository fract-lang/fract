/*
	ToInt8 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToInt8 String to 8bit integer.
// value Value to parse.
func ToInt8(value string) (int8, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(result), err
}
