/*
	ApplyEmbedFunctions Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// ApplyEmbedFunctions Add embed functions to interpreter source.
func (i *Interpreter) ApplyEmbedFunctions() {
	i.functions = append(i.functions,
		obj.Function{ // input function.
			Name:                  "input",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []obj.Parameter{
				{
					Name: "input.message",
					Default: obj.Value{
						Content: []string{""},
						Type:    fract.VALString,
					},
				},
			},
		},
	)
}
