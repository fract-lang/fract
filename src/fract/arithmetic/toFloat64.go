/*
	ToFloat64 Function.
*/

package arithmetic

import (
	"strconv"

	"github.com/fract-lang/src/grammar"
)

// ToFloat64 String to double.
// value Value to parse.
func ToFloat64(value string) (float64, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}
