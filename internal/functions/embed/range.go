package embed

import (
	"fmt"
	"strconv"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

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
