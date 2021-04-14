package functions

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// String Convert object to string.
// f Function.
// parameters Parameters.
func String(f obj.Function, parameters []obj.Variable) obj.Value {
	switch parameters[1].Value.Content[0].Data {
	case "parse":
		str := ""

		if value := parameters[0].Value; value.Array {
			if len(value.Content) == 0 {
				str = "[]"
			} else {
				sb := strings.Builder{}
				sb.WriteRune('[')
				for _, data := range value.Content {
					sb.WriteString(data.Data + " ")
				}
				str = sb.String()[:sb.Len()-1] + "]"
			}
		} else {
			str = parameters[0].Value.Content[0].Data
		}

		return obj.Value{
			Content: []obj.DataFrame{
				{
					Data: str,
					Type: fract.VALString,
				},
			},
		}
	case "object":
		fallthrough
	default: // Objects.
		return obj.Value{
			Content: []obj.DataFrame{
				{
					Data: fmt.Sprint(parameters[0].Value),
					Type: fract.VALString,
				},
			},
		}
	}
}
