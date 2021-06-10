package objects

import (
	"fmt"
	"strings"
)

// Value intance.
type Value struct {
	Content []Data
	Array   bool
}

func (v *Value) Print() bool {
	if v.Content == nil {
		return false
	}

	if v.Array {
		if len(v.Content) == 0 {
			fmt.Print("[]")
		} else {
			var sb strings.Builder
			sb.WriteByte('[')
			for _, data := range v.Content {
				sb.WriteString(data.Format() + " ")
			}
			fmt.Print(sb.String()[:sb.Len()-1] + "]")
		}
	} else {
		fmt.Print(v.Content[0].Format())
	}
	return true
}
