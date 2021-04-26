package functions

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Float Convert object to float.
// f Function.
// parameters Parameters.
func Float(f obj.Function, parameters []*obj.Variable) obj.Value {
	return obj.Value{
		Content: []obj.DataFrame{
			{
				Data: fmt.Sprintf(fract.FloatFormat, arithmetic.ToArithmetic(parameters[0].Value.Content[0].Data)),
				Type: fract.VALFloat,
			},
		},
	}
}
