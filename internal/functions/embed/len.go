package embed

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/objects"
)

// Len returns length of object.
func Len(f objects.Function, parameters []objects.Variable) objects.Value {
	parameter := parameters[0].Value
	if parameter.Array {
		return objects.Value{
			Content: []objects.Data{{Data: fmt.Sprint(len(parameter.Content))}},
		}
	} else if parameter.Content[0].Type == objects.VALString {
		return objects.Value{
			Content: []objects.Data{{Data: fmt.Sprint(len(parameter.Content[0].String()))}},
		}
	}
	return objects.Value{Content: []objects.Data{{Data: "0"}}}
}
