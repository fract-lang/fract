package make

import (
	"os"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
)

// Check invalid state of value.
func Check(value string) bool {
	if strings.HasSuffix(value, fract.FractExtension) {
		return true
	}
	value += ".fract"
	info, err := os.Stat(value)
	return err == nil && !info.IsDir()
}
