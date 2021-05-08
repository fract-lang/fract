package functions

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Float convert object to float.
func Float(f objects.Function, parameters []objects.Variable) objects.Value {
	return objects.Value{
		Content: []objects.DataFrame{
			{
				Data: fmt.Sprintf(fract.FloatFormat, arithmetic.ToArithmetic(parameters[0].Value.Content[0].Data)),
				Type: fract.VALFloat,
			},
		},
	}
}
