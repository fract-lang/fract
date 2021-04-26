package functions

import (
	"fmt"
	"strconv"

	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Range of object.
// f Function.
// parameters Parameters.
func Range(f obj.Function, parameters []*obj.Variable) obj.Value {
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

	if start.Content[0].Type != fract.VALInteger &&
		start.Content[0].Type != fract.VALFloat ||
		to.Content[0].Type != fract.VALInteger &&
			to.Content[0].Type != fract.VALFloat ||
		step.Content[0].Type != fract.VALInteger &&
			step.Content[0].Type != fract.VALFloat {
		fract.Error(f.Tokens[0][0], "Values should be integer or float!")
	}

	startV, _ := strconv.ParseFloat(start.Content[0].Data, 64)
	toV, _ := strconv.ParseFloat(to.Content[0].Data, 64)
	stepV, _ := strconv.ParseFloat(step.Content[0].Data, 64)

	if stepV <= 0 {
		return obj.Value{
			Content: nil,
			Array:   true,
		}
	}

	var dtype int16
	if start.Content[0].Type == fract.VALFloat ||
		to.Content[0].Type == fract.VALFloat ||
		step.Content[0].Type == fract.VALFloat {
		dtype = fract.VALFloat
	}

	returnValue := obj.Value{
		Content: []obj.DataFrame{},
		Array:   true,
	}

	if startV <= toV {
		for ; startV <= toV; startV += stepV {
			data := obj.DataFrame{
				Data: fmt.Sprintf(fract.FloatFormat, startV),
				Type: dtype,
			}
			data.Data = fract.FormatData(data)
			returnValue.Content = append(returnValue.Content, data)
		}
	} else {
		for ; startV >= toV; startV -= stepV {
			data := obj.DataFrame{
				Data: fmt.Sprintf(fract.FloatFormat, startV),
				Type: dtype,
			}
			data.Data = fract.FormatData(data)
			returnValue.Content = append(returnValue.Content, data)
		}
	}

	return returnValue
}
