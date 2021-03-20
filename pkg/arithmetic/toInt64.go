/*
	ToInt64 Function.
*/

package arithmetic

import (
	"strconv"
)

// ToInt64 String to 64bit integer.
// value Value to parse.
func ToInt64(value string) (int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(result), err
}
