package embed

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// String convert object to string.
func String(f objects.Function, parameters []objects.Variable) objects.Value {
	switch parameters[1].Value.Content[0].Data {
	case "parse":
		str := ""
		if value := parameters[0].Value; value.Array {
			if len(value.Content) == 0 {
				str = "[]"
			} else {
				var sb strings.Builder
				sb.WriteByte('[')
				for _, data := range value.Content {
					sb.WriteString(data.Data + " ")
				}
				str = sb.String()[:sb.Len()-1] + "]"
			}
		} else {
			str = parameters[0].Value.Content[0].Data
		}
		return objects.Value{
			Content: []objects.DataFrame{
				{
					Data: str,
					Type: fract.VALString,
				},
			},
		}
	case "bytecode":
		value := parameters[0].Value
		var sb strings.Builder
		for _, data := range value.Content {
			if data.Type != fract.VALInteger {
				sb.WriteByte(' ')
			}
			result, _ := strconv.ParseInt(data.Data, 10, 32)
			sb.WriteByte(byte(result))
		}
		return objects.Value{
			Content: []objects.DataFrame{
				{
					Data: sb.String(),
					Type: fract.VALString,
				},
			},
		}
	default: // Object.
		return objects.Value{
			Content: []objects.DataFrame{
				{
					Data: fmt.Sprint(parameters[0].Value),
					Type: fract.VALString,
				},
			},
		}
	}
}
