/*
	ToInt Function.
*/

package arithmetic

import (
	"strconv"

	"github.com/fract-lang/fract/pkg/grammar"
)

// ToInt String to integer.
// value Value to parse.
func ToInt(value string) (int, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, err
}
