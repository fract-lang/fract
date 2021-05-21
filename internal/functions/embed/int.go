package embed

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Int convert object to integer.
func Int(f objects.Function, parameters []objects.Variable) objects.Value {
	switch parameters[1].Value.Content[0].Data { // Cast type.
	case "strcode":
		var value objects.Value
		for _, byt := range []byte(parameters[0].Value.Content[0].Data) {
			value.Content = append(value.Content,
				objects.DataFrame{
					Data: fmt.Sprint(byt),
					Type: fract.VALInteger,
				})
		}
		value.Array = len(value.Content) > 1
		return value
	default: // Object.
		return objects.Value{
			Content: []objects.DataFrame{
				{
					Data: fmt.Sprint(int(arithmetic.ToArithmetic(parameters[0].Value.Content[0].Data))),
					Type: fract.VALInteger,
				},
			},
		}
	}
}
