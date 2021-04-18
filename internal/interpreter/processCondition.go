/*
	processCondition Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

var (
	// TrueValueIns True condition value instance.
	TrueValueIns obj.Value = obj.Value{
		Array:   false,
		Content: []obj.DataFrame{{Data: grammar.KwTrue}},
	}
)

// compare Compare values by operator.
// value0 First value of comparison.
// value1 Second value of comparison.
// operator Operator of comparison.
func compare(value0, value1 obj.Value, operator string) bool {
	// compare_values Compare values by operator.
	// data0 First data to compare.
	// data1 Secondary data to compare.
	compare_values := func(data0 obj.DataFrame, data1 obj.DataFrame) bool {
		if data0.Type != data1.Type &&
			(data0.Type == fract.VALString || data1.Type == fract.VALString) {
			return false
		}

		switch operator {
		case grammar.Equals: // Equals.
			if data0.Data != data1.Data {
				return false
			}
		case grammar.NotEquals: // Not equals.
			if data0.Data == data1.Data {
				return false
			}
		case grammar.TokenGreat: // Greater.
			if (data0.Type == fract.VALString && data0.Data <= data1.Data) ||
				(data0.Type != fract.VALString &&
					arithmetic.ToArithmetic(data0.Data) <=
						arithmetic.ToArithmetic(data1.Data)) {
				return false
			}
		case grammar.TokenLess: // Less.
			if (data0.Type == fract.VALString && data0.Data >= data1.Data) ||
				(data0.Type != fract.VALString &&
					arithmetic.ToArithmetic(data0.Data) >=
						arithmetic.ToArithmetic(data1.Data)) {
				return false
			}
		case grammar.GreaterEquals: // Greater or equals.
			if (data0.Type == fract.VALString && data0.Data < data1.Data) ||
				(data0.Type != fract.VALString &&
					arithmetic.ToArithmetic(data0.Data) <
						arithmetic.ToArithmetic(data1.Data)) {
				return false
			}
		case grammar.LessEquals: // Less or equals.
			if (data0.Type == fract.VALString && data0.Data > data1.Data) ||
				(data0.Type != fract.VALString &&
					arithmetic.ToArithmetic(data0.Data) >
						arithmetic.ToArithmetic(data1.Data)) {
				return false
			}
		}

		return true
	}

	// String comparison.
	if !value0.Array || !value1.Array {
		data0 := value0.Content[0]
		data1 := value1.Content[0]
		if (data0.Type == fract.VALString && data1.Type != fract.VALString) ||
			(data0.Type != fract.VALString && data1.Type == fract.VALString) {
			return false
		}

		return compare_values(data0, data1)
	}

	// Array comparison.
	if value0.Array || value1.Array {
		if (value0.Array && !value1.Array) ||
			(!value0.Array && value1.Array) {
			return false
		}

		if len(value0.Content) != len(value1.Content) {
			return operator == grammar.NotEquals
		}

		for index := range value0.Content {
			if !compare_values(value0.Content[index], value1.Content[index]) {
				return false
			}
		}

		return true
	}

	// Single value comparison.
	return compare_values(value0.Content[0], value1.Content[0])
}

// processCondition Process conditional expression and return result.
// tokens Tokens to process.
func (i *Interpreter) processCondition(tokens *[]obj.Token) string {
	i.processRange(tokens)

	// Process condition.
	ors := parser.DecomposeConditionalProcess(*tokens, grammar.LogicalOr)
	for _, current := range *ors {
		// Decompose and conditions.
		ands := parser.DecomposeConditionalProcess(current, grammar.LogicalAnd)
		// Is and long statement?
		if len(*ands) > 1 {
			for aindex := range *ands {
				if !compare(i.processValue(
					&(*ands)[aindex]), TrueValueIns, grammar.Equals) {
					return grammar.KwFalse
				}
			}
			return grammar.KwTrue
		}

		operatorIndex, operator := parser.FindConditionOperator(current)

		// Operator is not found?
		if operatorIndex == -1 {
			if compare(i.processValue(&current), TrueValueIns, grammar.Equals) {
				return grammar.KwTrue
			}
			continue
		}

		// Operator is first or last?
		if operatorIndex == 0 {
			fract.Error(current[0], "Comparison values are missing!")
		} else if operatorIndex == len(current)-1 {
			fract.Error(current[len(current)-1], "Comparison values are missing!")
		}

		if compare(i.processValue(
			vector.Sublist(current, 0, operatorIndex)), i.processValue(
			vector.Sublist(current, operatorIndex+1,
				len(current)-operatorIndex-1)), operator) {
			return grammar.KwTrue
		}
	}

	return grammar.KwFalse
}
