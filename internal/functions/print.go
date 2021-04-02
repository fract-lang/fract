package functions

import (
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Print values to CLI.
// f Function.
// parameters Parameters.
func Print(f obj.Function, parameters []obj.Value) {
	fract.PrintValue(parameters[0])
	fract.PrintValue(parameters[1])
}
