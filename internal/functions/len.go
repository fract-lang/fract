package functions

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Len of object.
// f Function.
// parameters Parameters.
func Len(f obj.Function, parameters []obj.Value) obj.Value {
	parameter := parameters[0]

	if parameter.Array {
		return obj.Value{
			Content: []string{fmt.Sprintf("%d", len(parameter.Content))},
		}
	}

	if parameter.Type == fract.VALString {
		return obj.Value{
			Content: []string{fmt.Sprintf("%d", len(parameter.Content[0]))},
		}
	}

	return obj.Value{
		Content: []string{"0"},
	}
}
