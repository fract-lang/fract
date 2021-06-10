package embed

import (
	"fmt"
	"strconv"
	"strings"

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
					sb.WriteString(data.String() + " ")
				}
				str = sb.String()[:sb.Len()-1] + "]"
			}
		} else {
			str = parameters[0].Value.Content[0].String()
		}
		return objects.Value{
			Content: []objects.Data{
				{
					Data: str,
					Type: objects.VALString,
				},
			},
		}
	case "bytecode":
		value := parameters[0].Value
		var sb strings.Builder
		for _, data := range value.Content {
			if data.Type != objects.VALInteger {
				sb.WriteByte(' ')
			}
			result, _ := strconv.ParseInt(data.String(), 10, 32)
			sb.WriteByte(byte(result))
		}
		return objects.Value{
			Content: []objects.Data{
				{
					Data: sb.String(),
					Type: objects.VALString,
				},
			},
		}
	default: // Object.
		return objects.Value{
			Content: []objects.Data{
				{
					Data: fmt.Sprint(parameters[0].Value),
					Type: objects.VALString,
				},
			},
		}
	}
}
