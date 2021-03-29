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
					Name: "message",
					Default: obj.Value{
						Content: []string{""},
						Type:    fract.VALString,
					},
				},
			},
		},
		obj.Function{ // exit function.
			Name:                  "exit",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []obj.Parameter{
				{
					Name: "code",
					Default: obj.Value{
						Content: []string{"0"},
						Type:    fract.VALInteger,
					},
				},
			},
		},
		obj.Function{ // len function.
			Name:                  "len",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: []obj.Parameter{
				{
					Name: "object",
				},
			},
		},
		obj.Function{ // range function.
			Name:                  "range",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []obj.Parameter{
				{
					Name: "start",
				},
				{
					Name: "to",
				},
				{
					Name: "step",
					Default: obj.Value{
						Content: []string{"1"},
					},
				},
			},
		},
	)
}
