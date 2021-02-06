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

// compare Compare values by operator.
// value0 First value of comparison.
// value1 Second value of comparison.
// operator Operator of comparison.
func compare(value0 float64, value1 float64, operator string) bool {
	switch operator {
	case grammar.TokenEquals: // Equals.
		return value0 == value1
	case grammar.NotEquals: // Not equals.
		return value0 != value1
	case grammar.TokenGreat: // Greater.
		return value0 > value1
	case grammar.TokenLess: // Less.
		return value0 < value1
	case grammar.GreaterEquals: // Greater or equals.
		return value0 >= value1
	default:
		return false
	}
}

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
	ors := parser.DecomposeConditionalProcess(tokens, grammar.TokenVerticalBar)
	for index := 0; index < ors.Len(); index++ {
		current := ors.At(index).(vector.Vector)

		// Decompose and conditions.
		ands := parser.DecomposeConditionalProcess(&current, grammar.TokenAmper)
		// Is and statement?
		if ands.Len() > 1 {
			for aindex := 0; aindex < ands.Len(); aindex++ {
				acurrent := ands.At(aindex).(vector.Vector)
				value, _ := arithmetic.ToFloat64(i.processValue(&acurrent).Content)
				if !compare(value, 1, grammar.TokenEquals) {
					return grammar.FALSE
				}
			}
			return grammar.TRUE
		}

		operatorIndex, operator := parser.FindConditionOperator(&current)

		// Operator is not found?
		if operatorIndex == -1 {
			value, _ := arithmetic.ToFloat64(i.processValue(&current).Content)
			if compare(value, 1, grammar.TokenEquals) {
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
		if compare(val0, val1, operator) {
			return grammar.TRUE
		}
	}

	return grammar.FALSE
}
