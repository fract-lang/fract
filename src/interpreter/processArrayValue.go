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
	value.Type = fract.VTInteger

	first := tokens.Vals[0].(objects.Token)

	// Initializer?
	if first.Value == grammar.TokenLBracket {
		valueList := tokens.Sublist(1, len(tokens.Vals)-2)
		if len(valueList.Vals) == 0 {
			fract.Error(first, "Size is not defined!")
		}

		val, _ := arithmetic.ToInt64(i.processValue(valueList).Content[0])
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
			value.Content = append(value.Content, i.processValue(
				tokens.Sublist(comma, index-comma)).Content...)
			comma = index + 1
		}
	}

	if comma < len(tokens.Vals)-1 {
		value.Content = append(value.Content, i.processValue(
			tokens.Sublist(comma, len(tokens.Vals)-comma-1)).Content...)
	}

	/* Set type to float if... */
	for index := range value.Content {
		current := value.Content[index]
		if strings.Index(current, grammar.TokenDot) != -1 ||
			strings.Index(current, grammar.TokenDot) != -1 {
			value.Type = fract.VTFloat
			break
		}
	}

	return value
}
