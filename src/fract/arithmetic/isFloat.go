/*
	IsFloat Function.
*/

package arithmetic

import (
	"regexp"
)

// IsFloat Value is an float?
// value Value to check.
func IsFloat(value string) bool {
	state, _ := regexp.MatchString("^(-|)\\s*[0-9]+(\\.[0-9]+)?$", value)
	return state
}
