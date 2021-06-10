package embed

import (
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/objects"
)

// Input returns input from command-line.
func Input(f objects.Function, parameters []objects.Variable) objects.Value {
	parameters[0].Value.Print()
	return objects.Value{
		Content: []objects.Data{
			{
				Data: cli.Input(""),
				Type: objects.VALString,
			},
		},
	}
}
