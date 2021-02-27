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

var (
	// TrueValueIns True condition value instance.
	TrueValueIns objects.Value = objects.Value{
		Array:   false,
		Type:    fract.VTInteger,
		Content: []string{"1"},
	}
)

// compare Compare values by operator.
// value0 First value of comparison.
// value1 Second value of comparison.
// operator Operator of comparison.
func compare(value0 objects.Value, value1 objects.Value, operator string) bool {
	if value0.Array != value1.Array || len(value0.Content) != len(value1.Content) {
		return false
	}
	for index := range value0.Content {
		item0, _ := arithmetic.ToFloat64(value0.Content[index])
		item1, _ := arithmetic.ToFloat64(value1.Content[index])
		switch operator {
		case grammar.TokenEquals: // Equals.
			if item0 == item1 {
			} else {
				return false
			}
		case grammar.NotEquals: // Not equals.
			if item0 != item1 {
			} else {
				return false
			}
		case grammar.TokenGreat: // Greater.
			if item0 > item1 {
			} else {
				return false
			}
		case grammar.TokenLess: // Less.
			if item0 < item1 {
			} else {
				return false
			}
		case grammar.GreaterEquals: // Greater or equals.
			if item0 >= item1 {
			} else {
				return false
			}
		case grammar.LessEquals: // Less or equals.
			if item0 <= item1 {
			} else {
				return false
			}
		}
	}
	return true
}

// processCondition Process conditional expression and return result.
// tokens Tokens to process.
func (i *Interpreter) processCondition(tokens *vector.Vector) int {
	i.processRange(tokens)

	// Process condition.
	ors := parser.DecomposeConditionalProcess(tokens, grammar.TokenVerticalBar)
	for index := range ors.Vals {
		current := ors.Vals[index].(*vector.Vector)

		// Decompose and conditions.
		ands := parser.DecomposeConditionalProcess(current, grammar.TokenAmper)
		// Is and long statement?
		if len(ands.Vals) > 1 {
			for aindex := range ands.Vals {
				if !compare(i.processValue(
					ands.Vals[aindex].(*vector.Vector)), TrueValueIns, grammar.TokenEquals) {
					return grammar.FALSE
				}
			}
			return grammar.TRUE
		}

		operatorIndex, operator := parser.FindConditionOperator(current)

		// Operator is not found?
		if operatorIndex == -1 {
			if compare(i.processValue(current), TrueValueIns, grammar.TokenEquals) {
				return grammar.TRUE
			}
			continue
		}

		// Operator is first or last?
		if operatorIndex == 0 {
			fract.Error(current.Vals[0].(objects.Token),
				"Comparison values are missing!")
		} else if operatorIndex == len(current.Vals)-1 {
			fract.Error(current.Vals[len(current.Vals)-1].(objects.Token),
				"Comparison values are missing!")
		}

		if compare(i.processValue(
			current.Sublist(0, operatorIndex)), i.processValue(
			current.Sublist(operatorIndex+1,
				len(current.Vals)-operatorIndex-1)), operator) {
			return grammar.TRUE
		}
	}

	return grammar.FALSE
}
