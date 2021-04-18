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

	formatData := func(data obj.DataFrame) string {
		if data.Type == VALFloat {
			for index := len(data.Data) - 1; index >= 0; index-- {
				if ch := data.Data[index]; ch != '0' {
					data.Data = data.Data[:index+1]
					if ch == '.' {
						data.Data += "0"
					}
					return data.Data
				}
			}
			return data.Data
		}
		return data.Data
	}

	if value.Array {
		if len(value.Content) == 0 {
			fmt.Print("[]")
		} else {
			sb := strings.Builder{}
			sb.WriteRune('[')
			for _, data := range value.Content {
				sb.WriteString(formatData(data) + " ")
			}
			fmt.Print(sb.String()[:sb.Len()-1] + "]")
		}
	} else {
		fmt.Print(formatData(value.Content[0]))
	}
	return true
}
