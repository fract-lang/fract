/*
	ToFloat64 Function.
*/

package arithmetic

import (
	"strconv"
)

// ToFloat64 String to double.
// value Value to parse.
func ToFloat64(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}
