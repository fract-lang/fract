/*
	processArrayValue Function
*/

package interpreter

import (
	"strings"

	"../fract"
	"../fract/arithmetic"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// processArrayValue Process array value.
// tokens Tokens.
func (i *Interpreter) processArrayValue(tokens *vector.Vector) objects.Value {
	var value objects.Value
	value.Array = true
	value.Type = fract.VTInteger

	first := tokens.Vals[0].(objects.Token)

	// Initializer?
	if first.Value == grammar.TokenLBracket {
		valueList := tokens.Sublist(1, len(tokens.Vals)-2)

		if len(valueList.Vals) == 0 {
			fract.Error(first, "Size is not defined!")
		}

		value := i.processValue(valueList)
		if value.Array {
			fract.Error(first, "Arrays is not used in array constructors!")
		}
		if value.Type == fract.VTFloat {
			fract.Error(first, "Float values is not used in array constructors!")
		}

		val, _ := arithmetic.ToInt64(value.Content[0])
		if val < 0 {
			fract.Error(first, "Value is not lower than zero!")
		}
		value.Content = make([]string, val)
		for index := range value.Content {
			value.Content[index] = "0"
		}
		return value
	}

	comma := 1
	for index := 1; index < len(tokens.Vals)-1; index++ {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeComma {
			val := i.processValue(tokens.Sublist(comma, index-comma))
			value.Content = append(value.Content, val.Content...)
			if !value.Charray {
				value.Charray = val.Charray
			}
			comma = index + 1
		}
	}

	if comma < len(tokens.Vals)-1 {
		val := i.processValue(tokens.Sublist(comma, len(tokens.Vals)-comma-1))
		value.Content = append(value.Content, val.Content...)
		if !value.Charray {
			value.Charray = val.Charray
		}
	}

	/* Set type to float if... */
	for index := range value.Content {
		current := value.Content[index]
		if strings.Index(current, grammar.TokenDot) != -1 {
			value.Type = fract.VTFloat
			break
		}
	}

	return value
}
