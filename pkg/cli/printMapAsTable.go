/*
	PrintMapAsTable Function.
*/

package cli

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/str"
)

// PrintMapAsTable Print map to cli screen as table.
// dict Map to print.
func PrintMapAsTable(dict map[string]string) {
	maxlen := 0
	for key := range dict {
		if maxlen < len(key) {
			maxlen = len(key)
		}
	}
	maxlen += 5
	for key := range dict {
		fmt.Println(key + " " + str.GetWhitespace(maxlen-len(key)) + dict[key])
	}
}
