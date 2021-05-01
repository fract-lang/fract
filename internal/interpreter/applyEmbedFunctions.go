package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

//! Embed functions should have a lowercase names.

// ApplyEmbedFunctions to interpreter source.
func (i *Interpreter) ApplyEmbedFunctions() {
	i.functions = append(i.functions,
		objects.Function{ // print function.
			Name:                  "print",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 2,
			Parameters: &[]objects.Parameter{
				{
					Name: "value",
					Default: objects.Value{
						Content: []objects.DataFrame{
							{
								Type: fract.VALString,
							},
						},
					},
				},
				{
					Name: "fin",
					Default: objects.Value{
						Content: []objects.DataFrame{
							{
								Data: "\n",
								Type: fract.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // input function.
			Name:                  "input",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: &[]objects.Parameter{
				{
					Name: "message",
					Default: objects.Value{
						Content: []objects.DataFrame{
							{
								Data: "",
								Type: fract.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // exit function.
			Name:                  "exit",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: &[]objects.Parameter{
				{
					Name: "code",
					Default: objects.Value{
						Content: []objects.DataFrame{{Data: "0"}},
					},
				},
			},
		},
		objects.Function{ // len function.
			Name:                  "len",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: &[]objects.Parameter{
				{
					Name: "object",
				},
			},
		},
		objects.Function{ // range function.
			Name:                  "range",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: &[]objects.Parameter{
				{
					Name: "start",
				},
				{
					Name: "to",
				},
				{
					Name: "step",
					Default: objects.Value{
						Content: []objects.DataFrame{{Data: "1"}},
					},
				},
			},
		},
		objects.Function{ // make function.
			Name:                  "make",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: &[]objects.Parameter{
				{
					Name: "size",
				},
			},
		},
		objects.Function{ // string function.
			Name:                  "string",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: &[]objects.Parameter{
				{
					Name: "object",
				},
				{
					Name: "type",
					Default: objects.Value{
						Content: []objects.DataFrame{
							{
								Data: "parse",
								Type: fract.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // int function.
			Name:                  "int",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: &[]objects.Parameter{
				{
					Name: "object",
				},
				{
					Name: "type",
					Default: objects.Value{
						Content: []objects.DataFrame{
							{
								Data: "parse",
								Type: fract.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // float function.
			Name:                  "float",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: &[]objects.Parameter{
				{
					Name: "object",
				},
			},
		},
	)
}
