/*
	PrintValue Function.
*/

package fract

import (
	"fmt"
	"math/big"
	"strings"

	obj "github.com/fract-lang/fract/pkg/objects"
)

func format(data obj.DataFrame) string {
	if data.Type == VALInteger {
		bigfloat, _ := new(big.Float).SetString(data.Data)
		data.Data = bigfloat.String()
	}

	return data.Data
}

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
			var sb strings.Builder
			sb.WriteByte('[')
			for _, data := range value.Content {
				sb.WriteString(format(data) + " ")
			}
			fmt.Print(sb.String()[:sb.Len()-1] + "]")
		}
	} else {
		fmt.Print(format(value.Content[0]))
	}
	return true
}
