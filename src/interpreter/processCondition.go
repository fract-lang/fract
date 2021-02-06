/*
	processCondition Function.
*/

package interpreter

import (
	"../fract"
	"../fract/arithmetic"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processCondition Process conditional expression and return result.
// tokens Tokens to process.
func (i *Interpreter) processCondition(tokens *vector.Vector) int {
	/* Check parentheses range. */
	for true {
		_range, found := parser.DecomposeBrace(tokens)

		/* Parentheses are not found! */
		if found == -1 {
			break
		}

		var _token objects.Token
		_token.Value = i.processValue(&_range).Content
		_token.Type = fract.TypeValue
		tokens.Insert(found, _token)
	}

	// Process condition.
	ors := parser.DecomposeOrConditionalProcess(tokens)
	for index := 0; index < ors.Len(); index++ {
		current := ors.At(index).(vector.Vector)
		operatorIndex, operator := findConditionOperator(&current)

		// Operator is not found?
		if operatorIndex == -1 {
			value, _ := arithmetic.ToFloat64(current.First().(objects.Token).Value)
			if value == 1 {
				return grammar.TRUE
			}
			continue
		}

		// Operator is first or last?
		if operatorIndex == 0 {
			fract.Error(current.First().(objects.Token), "Comparison values are missing!")
		} else if operatorIndex == current.Len()-1 {
			fract.Error(current.Last().(objects.Token), "Comparison values are missing!")
		}

		val0L := current.Sublist(0, operatorIndex)
		val1L := current.Sublist(operatorIndex+1, current.Len()-operatorIndex-1)
		val0, _ := arithmetic.ToFloat64(i.processValue(&val0L).Content)
		val1, _ := arithmetic.ToFloat64(i.processValue(&val1L).Content)
		switch operator {
		case grammar.TokenEquals:
			if val0 == val1 {
				return grammar.TRUE
			}
		}
	}
	return grammar.FALSE
}
