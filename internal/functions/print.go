package functions

import (
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Print values to CLI.
// f Function.
// parameters Parameters.
func Print(f obj.Function, parameters []obj.Variable) {
	if parameters[0].Value.Content == nil {
		fract.Error(f.Tokens[0][0], "Invalid value!")
	}
	fract.PrintValue(parameters[0].Value)
	fract.PrintValue(parameters[1].Value)
}
