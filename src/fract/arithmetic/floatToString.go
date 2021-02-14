/*
	FloatToString Function.
*/

package arithmetic

import (
	"fmt"
)

// FloatToString Float to string.
// value Value to parse.
func FloatToString(value interface{}) string {
	return fmt.Sprintf("%f", value)
}
