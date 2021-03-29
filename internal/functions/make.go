package functions

import (
	"strconv"

	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Make array by size.
// f Function.
// parameters Parameters.
func Make(f obj.Function, parameters []obj.Value) obj.Value {
	size := parameters[0]

	if size.Array {
		fract.Error(f.Tokens[0][0], "Array is not a valid value!")
	} else if size.Type != fract.VALInteger {
		fract.Error(f.Tokens[0][0], "Exit code is only be integer!")
	}

	sizev, _ := strconv.ParseInt(size.Content[0], 10, 64)
	return obj.Value{
		Content: make([]string, sizev),
		Array:   true,
	}
}
