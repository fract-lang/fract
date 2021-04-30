package functions

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Int Convert object to integer.
// f Function.
// parameters Parameters.
func Int(f obj.Function, parameters []*obj.Variable) obj.Value {
	switch parameters[1].Value.Content[0].Data { // Cast type.
	case "strcode":
		var value obj.Value
		for _, byt := range []byte(parameters[0].Value.Content[0].Data) {
			value.Content = append(value.Content,
				obj.DataFrame{
					Data: fmt.Sprint(byt),
					Type: fract.VALInteger,
				})
		}
		value.Array = len(value.Content) > 1
		return value
	default: // Object.
		return obj.Value{
			Content: []obj.DataFrame{
				{
					Data: fmt.Sprint(int(arithmetic.ToArithmetic(parameters[0].Value.Content[0].Data))),
					Type: fract.VALInteger,
				},
			},
		}
	}
}
