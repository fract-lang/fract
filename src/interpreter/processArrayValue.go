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
	value.Type = fract.VTIntegerArray

	first := tokens.Vals[0].(objects.Token)

	// Initializer?
	if first.Value == grammar.TokenLBracket {
		last := tokens.Vals[len(tokens.Vals)-1].(objects.Token)
		if last.Type != fract.TypeBrace && last.Value != grammar.TokenRBracket {
			fract.Error(last, "Array close bracket is not found!")
		}
		valueList := tokens.Sublist(1, len(tokens.Vals)-2)
		val, err := arithmetic.ToInt64(i.processValue(&valueList).Content[0])
		if err != nil {
			fract.Error(first, "Value out of range!")
		}
		if val < 0 {
			fract.Error(first, "Value is not lower than zero!")
		}
		value.Content = make([]string, val)
		for index := range value.Content {
			value.Content[index] = "0"
		}
		return value
	}

	if first.Type != fract.TypeBrace && first.Value != grammar.TokenLBrace {
		fract.Error(first, "Array brace is not found!")
	}
	last := tokens.Vals[len(tokens.Vals)-1].(objects.Token)
	if last.Type != fract.TypeBrace && last.Value != grammar.TokenRBrace {
		fract.Error(last, "Array close brace is not found!")
	}

	comma := 1
	for index := 1; index < len(tokens.Vals)-1; index++ {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeComma {
			valueList := tokens.Sublist(comma, index-comma)
			value.Content = append(value.Content, i.processValue(&valueList).Content...)
			comma = index + 1
		}
	}

	if comma < len(tokens.Vals)-1 {
		valueList := tokens.Sublist(comma, len(tokens.Vals)-comma-1)
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
