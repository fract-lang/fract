package functions

import (
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Input from cli.
// f Function.
// parameters Parameters.
func Input(f obj.Function, parameters []obj.Value) obj.Value {
	fract.PrintValue(parameters[0])
	return obj.Value{
		Content: []string{cli.Input("")},
		Type:    fract.VALString,
	}
}
