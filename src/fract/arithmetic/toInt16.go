/*
	ToInt16 Function.
*/

package arithmetic

import (
	"strconv"

	"../../grammar"
)

// ToInt16 String to 16bit integer.
// value Value to parse.
func ToInt16(value string) (int16, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(result), err
}
