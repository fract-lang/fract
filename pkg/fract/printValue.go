/*
	PrintValue Function.
*/

package fract

import (
	"fmt"
	"strings"

	obj "github.com/fract-lang/fract/pkg/objects"
)

// PrintValue Print value to screen.
// value Value to print.
func PrintValue(value obj.Value) bool {
	if value.Content == nil {
		return false
	}

	if value.Array {
		if len(value.Content) == 0 {
			fmt.Print("[]")
		} else {
			sb := strings.Builder{}
			sb.WriteByte('[')
			for _, data := range value.Content {
				sb.WriteString(FormatData(data) + " ")
			}
			fmt.Print(sb.String()[:sb.Len()-1] + "]")
		}
	} else {
		fmt.Print(FormatData(value.Content[0]))
	}
	return true
}
