/*
	PrintMapAsTable Function.
*/

package cli

import (
	"fmt"
)

// *********************
//       PRIVATE
// *********************

// Returns string whitespace by count.
// count Count of whitespace.
func getws(count int) string {
	var str string = ""
	for counter := 1; counter <= count; counter++ {
		str += " "
	}
	return str
}

// *********************
//        PUBLIC
// *********************

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
		fmt.Println(key + " " + getws(maxlen-len(key)) + dict[key])
	}
}
