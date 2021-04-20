/*
	Check Function.
*/

package make

import (
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/fs"
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
