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
	return obj.Value{
		Content: []obj.DataFrame{
			{
				Data: fmt.Sprint(int64(arithmetic.ToArithmetic(parameters[0].Value.Content[0].Data))),
				Type: fract.VALInteger,
			},
		},
	}
}
