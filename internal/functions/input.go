package functions

import (
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Input returns input from command-line.
func Input(f objects.Function, parameters []objects.Variable) objects.Value {
	fract.PrintValue(parameters[0].Value)
	return objects.Value{
		Content: []objects.DataFrame{
			{
				Data: cli.Input(""),
				Type: fract.VALString,
			},
		},
	}
}
