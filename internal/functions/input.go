package functions

import (
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Input from cli.
// f Function.
// parameters Parameters.
func Input(f obj.Function, parameters []*obj.Variable) obj.Value {
	fract.PrintValue(parameters[0].Value)
	return obj.Value{
		Content: []obj.DataFrame{
			{
				Data: cli.Input(""),
				Type: fract.VALString,
			},
		},
	}
}
