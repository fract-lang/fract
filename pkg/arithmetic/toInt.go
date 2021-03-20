/*
	ToInt Function.
*/

package arithmetic

import (
	"strconv"
)

// ToInt String to integer.
// value Value to parse.
func ToInt(value string) (int, error) {
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, err
}
