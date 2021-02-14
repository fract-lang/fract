/*
	IntToString Function.
*/

package arithmetic

import (
	"fmt"
)

// IntToString Integer to string.
// value Value to parse.
func IntToString(value interface{}) string {
	return fmt.Sprintf("%d", value)
}
