/*
	processArrayValue Function
*/

package interpreter

import (
	"strings"

	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// processArrayValue Process array value.
// tokens Tokens.
func (i *Interpreter) processArrayValue(tokens *vector.Vector) objects.Value {
	first := tokens.First().(objects.Token)
	if first.Type != fract.TypeBrace || first.Value != grammar.TokenLBracket {
		fract.Error(first, "Array bracket is not found!")
	}
	last := tokens.Last().(objects.Token)
	if last.Type != fract.TypeBrace || last.Value != grammar.TokenRBracket {
		fract.Error(last, "Array close bracket is not found!")
	}
	var value objects.Value
	value.Type = fract.VTIntegerArray

	comma := 1
	for index := 1; index < tokens.Len()-1; index++ {
		current := tokens.At(index).(objects.Token)
		if current.Type == fract.TypeComma {
			valueList := tokens.Sublist(comma, index-comma)
			value.Content = append(value.Content, i.processValue(&valueList).Content...)
			comma = index + 1
		}
	}

	if comma < tokens.Len()-1 {
		valueList := tokens.Sublist(comma, tokens.Len()-comma-1)
		value.Content = append(value.Content, i.processValue(&valueList).Content...)
	}

	/* Set type to float if... */
	for index := range value.Content {
		current := value.Content[index]
		if strings.Index(current, grammar.TokenDot) != -1 ||
			strings.Index(current, grammar.TokenDot) != -1 {
			value.Type = fract.VTFloatArray
			break
		}
	}

	return value
}
