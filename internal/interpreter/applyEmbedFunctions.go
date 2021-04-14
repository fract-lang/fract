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
		obj.Function{ // print function.
			Name:                  "print",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []obj.Parameter{
				{
					Name: "value",
				},
				{
					Name: "fin",
					Default: obj.Value{
						Content: []obj.DataFrame{
							{
								Data: "\n",
								Type: fract.VALString,
							},
						},
					},
				},
			},
		},
		obj.Function{ // input function.
			Name:                  "input",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []obj.Parameter{
				{
					Name: "message",
					Default: obj.Value{
						Content: []obj.DataFrame{
							{
								Data: "",
								Type: fract.VALString,
							},
						},
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
						Content: []obj.DataFrame{{Data: "0"}},
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
						Content: []obj.DataFrame{{Data: "1"}},
					},
				},
			},
		},
		obj.Function{ // make function.
			Name:                  "make",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: []obj.Parameter{
				{
					Name: "size",
				},
			},
		},
		obj.Function{ // string function.
			Name:                  "string",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []obj.Parameter{
				{
					Name: "object",
				},
				{
					Name: "type",
					Default: obj.Value{
						Content: []obj.DataFrame{
							{
								Data: "object",
								Type: fract.VALString,
							},
						},
					},
				},
			},
		},
		obj.Function{ // int function.
			Name:                  "int",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: []obj.Parameter{
				{
					Name: "object",
				},
			},
		},
		obj.Function{ // float function.
			Name:                  "float",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: []obj.Parameter{
				{
					Name: "object",
				},
			},
		},
	)
}
