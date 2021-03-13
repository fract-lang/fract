/*
	Check Function.
*/

package make

import (
	"strings"

	"../../../fract"
	"../../../utils/fs"
)

// Check Check invalid state of value.
// value Value to check.
func Check(value string) bool {
	if strings.HasSuffix(value, fract.FractExtension) {
		return true
	}
	value += ".fract"
	return fs.ExistFile(value)
}
