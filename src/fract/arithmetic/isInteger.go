/*
	IsInteger Function.
*/

package arithmetic

import (
	"regexp"
)

// IsInteger Value is an integer?
// value Value to check.
func IsInteger(value string) bool {
	state, _ := regexp.MatchString("^(-|)\\s*[0-9]+$", value)
	return state
}
