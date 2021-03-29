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
func Range(f obj.Function, parameters []obj.Value) obj.Value {
	start := parameters[0]
	to := parameters[1]
	step := parameters[2]

	if start.Array {
		fract.Error(f.Tokens[0][0], "'start' argument should be numeric!")
	} else if to.Array {
		fract.Error(f.Tokens[0][0], "'to' argument should be numeric!")
	} else if step.Array {
		fract.Error(f.Tokens[0][0], "'step' argument should be numeric!")
	}

	startV, _ := strconv.ParseFloat(start.Content[0], 64)
	toV, _ := strconv.ParseFloat(to.Content[0], 64)
	stepV, _ := strconv.ParseFloat(step.Content[0], 64)

	if startV > toV || stepV <= 0 {
		return obj.Value{
			Content: nil,
			Array:   true,
		}
	}

	returnValue := obj.Value{
		Content: []string{},
		Array:   true,
	}

	for ; startV <= toV; startV += stepV {
		returnValue.Content = append(returnValue.Content, fmt.Sprintf("%g", startV))
	}

	return returnValue
}
