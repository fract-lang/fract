package embed

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Len returns length of object.
func Len(f objects.Function, parameters []objects.Variable) objects.Value {
	parameter := parameters[0].Value
	if parameter.Array {
		return objects.Value{
			Content: []objects.DataFrame{{Data: fmt.Sprint(len(parameter.Content))}},
		}
	} else if parameter.Content[0].Type == fract.VALString {
		return objects.Value{
			Content: []objects.DataFrame{{Data: fmt.Sprint(len(parameter.Content[0].Data))}},
		}
	}
	return objects.Value{Content: []objects.DataFrame{{Data: "0"}}}
}
