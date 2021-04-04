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
	} else if size.Content[0].Type != fract.VALInteger {
		fract.Error(f.Tokens[0][0], "Exit code is only be integer!")
	}

	sizev, _ := strconv.ParseInt(size.Content[0].Data, 10, 64)
	if sizev < 0 {
		fract.Error(f.Tokens[0][0], "Size should be minimum zero!")
	}
	return obj.Value{
		Content: make([]obj.DataFrame, sizev),
		Array:   true,
	}
}
