package functions

import (
	"strconv"

	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Make array by size.
// f Function.
// parameters Parameters.
func Make(f obj.Function, parameters []*obj.Variable) obj.Value {
	size := parameters[0].Value

	if size.Array {
		fract.Error(f.Tokens[0][0], "Array is not a valid value!")
	} else if size.Content[0].Type != fract.VALInteger {
		fract.Error(f.Tokens[0][0], "Exit code is only be integer!")
	}

	sizev, _ := strconv.Atoi(size.Content[0].Data)
	if sizev < 0 {
		fract.Error(f.Tokens[0][0], "Size should be minimum zero!")
	}

	value := obj.Value{Array: true}

	if sizev > 0 {
		var index int
		for ; index < sizev; index++ {
			value.Content = append(value.Content, obj.DataFrame{Data: "0"})
		}
	} else {
		value.Content = []obj.DataFrame{}
	}

	return value
}
