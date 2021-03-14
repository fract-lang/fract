/*
	ToInt64 Function.
*/

package arithmetic

import (
	"strconv"

	"github.com/fract-lang/fract/src/grammar"
)

// ToInt64 String to 64bit integer.
// value Value to parse.
func ToInt64(value string) (int64, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(result), err
}
