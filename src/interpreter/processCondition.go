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

// checkEmpty Check value is empty. And return 0 if empty.
// values Values to check.
func checkEmpty(values []string) []string {
	if len(values) == 0 {
		return []string{"0"}
	}
	return values
}

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
	case grammar.LessEquals: // Less or equals.
		return value0 <= value1
	default:
		return false
	}
}

// processCondition Process conditional expression and return result.
// tokens Tokens to process.
func (i *Interpreter) processCondition(tokens *vector.Vector) int {
	/* Check parentheses range. */
	for true {
		_range, found := parser.DecomposeBrace(tokens, grammar.TokenLParenthes,
			grammar.TokenRParenthes)

		/* Parentheses are not found! */
		if found == -1 {
			break
		}

		var _token objects.Token
		_token.Value = i.processValue(_range).Content[0]
		_token.Type = fract.TypeValue
		tokens.Insert(found, _token)
	}
	// Process condition.
	ors := parser.DecomposeConditionalProcess(tokens, grammar.TokenVerticalBar)
	for index := range ors.Vals {
		current := ors.Vals[index].(vector.Vector)

		// Decompose and conditions.
		ands := parser.DecomposeConditionalProcess(&current, grammar.TokenAmper)
		// Is and statement?
		if len(ands.Vals) > 1 {
			for aindex := range ands.Vals {
				acurrent := ands.Vals[aindex].(vector.Vector)
				value, _ := arithmetic.ToFloat64(i.processValue(&acurrent).Content[0])
				if !compare(value, 1, grammar.TokenEquals) {
					return grammar.FALSE
				}
			}
			return grammar.TRUE
		}

		operatorIndex, operator := parser.FindConditionOperator(&current)

		// Operator is not found?
		if operatorIndex == -1 {
			value, _ := arithmetic.ToFloat64(i.processValue(&current).Content[0])
			if compare(value, 1, grammar.TokenEquals) {
				return grammar.TRUE
			}
			continue
		}

		// Operator is first or last?
		if operatorIndex == 0 {
			fract.Error(current.Vals[0].(objects.Token), "Comparison values are missing!")
		} else if operatorIndex == len(current.Vals)-1 {
			fract.Error(current.Vals[len(current.Vals)-1].(objects.Token), "Comparison values are missing!")
		}

		val0, _ := arithmetic.ToFloat64(checkEmpty(i.processValue(
			current.Sublist(0, operatorIndex)).Content)[0])
		val1, _ := arithmetic.ToFloat64(checkEmpty(i.processValue(
			current.Sublist(operatorIndex+1, len(current.Vals)-operatorIndex-1)).Content)[0])
		if compare(val0, val1, operator) {
			return grammar.TRUE
		}
	}

	return grammar.FALSE
}
