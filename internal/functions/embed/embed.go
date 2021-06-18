package embed

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/cli"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// Exit from application with code.
func Exit(f objects.Function, parameters []objects.Variable) {
	code := parameters[0].Value
	if code.Array {
		fract.Error(f.Tokens[0][0], "Array is not a valid value!")
	} else if code.Content[0].Type != objects.VALInteger {
		fract.Error(f.Tokens[0][0], "Exit code is only be integer!")
	}
	exit_code, _ := strconv.ParseInt(code.Content[0].String(), 10, 64)
	os.Exit(int(exit_code))
}

// Float convert object to float.
func Float(f objects.Function, parameters []objects.Variable) objects.Value {
	return objects.Value{
		Content: []objects.Data{
			{
				Data: fmt.Sprintf(fract.FloatFormat, arithmetic.ToArithmetic(parameters[0].Value.Content[0].String())),
				Type: objects.VALFloat,
			},
		},
	}
}

// Input returns input from command-line.
func Input(f objects.Function, parameters []objects.Variable) objects.Value {
	parameters[0].Value.Print()
	return objects.Value{
		Content: []objects.Data{
			{
				Data: cli.Input(""),
				Type: objects.VALString,
			},
		},
	}
}

// Int convert object to integer.
func Int(f objects.Function, parameters []objects.Variable) objects.Value {
	switch parameters[1].Value.Content[0].Data { // Cast type.
	case "strcode":
		var value objects.Value
		for _, byt := range []byte(parameters[0].Value.Content[0].String()) {
			value.Content = append(value.Content,
				objects.Data{
					Data: fmt.Sprint(byt),
					Type: objects.VALInteger,
				})
		}
		value.Array = len(value.Content) > 1
		return value
	default: // Object.
		return objects.Value{
			Content: []objects.Data{
				{
					Data: fmt.Sprint(int(arithmetic.ToArithmetic(parameters[0].Value.Content[0].String()))),
					Type: objects.VALInteger,
				},
			},
		}
	}
}

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

// Make array by size.
func Make(f objects.Function, parameters []objects.Variable) objects.Value {
	size := parameters[0].Value
	if size.Array {
		fract.Error(f.Tokens[0][0], "Array is not a valid value!")
	} else if size.Content[0].Type != objects.VALInteger {
		fract.Error(f.Tokens[0][0], "Exit code is only be integer!")
	}

	sizev, _ := strconv.Atoi(size.Content[0].String())
	if sizev < 0 {
		fract.Error(f.Tokens[0][0], "Size should be minimum zero!")
	}

	value := objects.Value{Array: true}
	if sizev > 0 {
		var index int
		for ; index < sizev; index++ {
			value.Content = append(value.Content, objects.Data{Data: "0"})
		}
	} else {
		value.Content = []objects.Data{}
	}
	return value
}

// Print values to cli.
func Print(f objects.Function, parameters []objects.Variable) {
	if parameters[0].Value.Content == nil {
		fract.Error(f.Tokens[0][0], "Invalid value!")
	}
	parameters[0].Value.Print()
	parameters[1].Value.Print()
}

// Range returns array by parameters.
func Range(f objects.Function, parameters []objects.Variable) objects.Value {
	start := parameters[0].Value
	to := parameters[1].Value
	step := parameters[2].Value
	if start.Array {
		fract.Error(f.Tokens[0][0], "'start' argument should be numeric!")
	} else if to.Array {
		fract.Error(f.Tokens[0][0], "'to' argument should be numeric!")
	} else if step.Array {
		fract.Error(f.Tokens[0][0], "'step' argument should be numeric!")
	}
	if start.Content[0].Type != objects.VALInteger &&
		start.Content[0].Type != objects.VALFloat || to.Content[0].Type != objects.VALInteger &&
		to.Content[0].Type != objects.VALFloat || step.Content[0].Type != objects.VALInteger &&
		step.Content[0].Type != objects.VALFloat {
		fract.Error(f.Tokens[0][0], "Values should be integer or float!")
	}

	startV, _ := strconv.ParseFloat(start.Content[0].String(), 64)
	toV, _ := strconv.ParseFloat(to.Content[0].String(), 64)
	stepV, _ := strconv.ParseFloat(step.Content[0].String(), 64)
	if stepV <= 0 {
		return objects.Value{
			Content: nil,
			Array:   true,
		}
	}

	var dtype uint8
	if start.Content[0].Type == objects.VALFloat || to.Content[0].Type == objects.VALFloat || step.Content[0].Type == objects.VALFloat {
		dtype = objects.VALFloat
	}
	returnValue := objects.Value{
		Content: []objects.Data{},
		Array:   true,
	}
	if startV <= toV {
		for ; startV <= toV; startV += stepV {
			data := objects.Data{
				Data: fmt.Sprintf(fract.FloatFormat, startV),
				Type: dtype,
			}
			data.Data = data.Format()
			returnValue.Content = append(returnValue.Content, data)
		}
	} else {
		for ; startV >= toV; startV -= stepV {
			data := objects.Data{
				Data: fmt.Sprintf(fract.FloatFormat, startV),
				Type: dtype,
			}
			data.Data = data.Format()
			returnValue.Content = append(returnValue.Content, data)
		}
	}
	return returnValue
}

// String convert object to string.
func String(f objects.Function, parameters []objects.Variable) objects.Value {
	switch parameters[1].Value.Content[0].Data {
	case "parse":
		str := ""
		if value := parameters[0].Value; value.Array {
			if len(value.Content) == 0 {
				str = "[]"
			} else {
				var sb strings.Builder
				sb.WriteByte('[')
				for _, data := range value.Content {
					sb.WriteString(data.String() + " ")
				}
				str = sb.String()[:sb.Len()-1] + "]"
			}
		} else {
			str = parameters[0].Value.Content[0].String()
		}
		return objects.Value{
			Content: []objects.Data{
				{
					Data: str,
					Type: objects.VALString,
				},
			},
		}
	case "bytecode":
		value := parameters[0].Value
		var sb strings.Builder
		for _, data := range value.Content {
			if data.Type != objects.VALInteger {
				sb.WriteByte(' ')
			}
			result, _ := strconv.ParseInt(data.String(), 10, 32)
			sb.WriteByte(byte(result))
		}
		return objects.Value{
			Content: []objects.Data{
				{
					Data: sb.String(),
					Type: objects.VALString,
				},
			},
		}
	default: // Object.
		return objects.Value{
			Content: []objects.Data{
				{
					Data: fmt.Sprint(parameters[0].Value),
					Type: objects.VALString,
				},
			},
		}
	}
}
