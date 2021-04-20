package functions

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Len of object.
// f Function.
// parameters Parameters.
func Len(f obj.Function, parameters []obj.Variable) obj.Value {
	parameter := parameters[0].Value

	if parameter.Array {
		return obj.Value{
			Content: []obj.DataFrame{{Data: fmt.Sprintf("%d", len(parameter.Content))}},
		}
	}

	if parameter.Content[0].Type == fract.VALString {
		return obj.Value{
			Content: []obj.DataFrame{{Data: fmt.Sprintf("%d", len(parameter.Content[0].Data))}},
		}
	}

	return obj.Value{Content: []obj.DataFrame{{Data: "0"}}}
}
