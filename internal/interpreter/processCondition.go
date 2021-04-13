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
	// String comparison.
	if !value0.Array || !value1.Array {
		data0 := value0.Content[0]
		data1 := value1.Content[0]
		if (data0.Type == fract.VALString && data1.Type != fract.VALString) ||
			(data0.Type != fract.VALString && data1.Type == fract.VALString) {
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
			if data0.Data <= data1.Data {
				return false
			}
		case grammar.TokenLess: // Less.
			if data0.Data >= data1.Data {
				return false
			}
		case grammar.GreaterEquals: // Greater or equals.
			if data0.Data < data1.Data {
				return false
			}
		case grammar.LessEquals: // Less or equals.
			if data0.Data > data1.Data {
				return false
			}
		}
		return true
	}

	// Array.
	if value0.Array || value1.Array {
		if (value0.Array && !value1.Array) ||
			(!value0.Array && value1.Array) {
			return false
		}

		get_sum := func(value obj.Value) float64 {
			sum := 0.
			for _, item := range value.Content {
				sum += arithmetic.ToArithmetic(item.Data)
			}
			return sum
		}

		switch operator {
		case grammar.Equals: // Equals.
			if len(value0.Content) != len(value1.Content) {
				return false
			}
		case grammar.NotEquals: // Not equals.
			if get_sum(value0) == get_sum(value1) {
				return false
			}
		case grammar.TokenGreat: // Greater.
			if get_sum(value0) <= get_sum(value1) {
				return false
			}
		case grammar.TokenLess: // Less.
			if get_sum(value0) >= get_sum(value1) {
				return false
			}
		case grammar.GreaterEquals: // Greater or equals.
			if get_sum(value0) < get_sum(value1) {
				return false
			}
		case grammar.LessEquals: // Less or equals.
			if get_sum(value0) > get_sum(value1) {
				return false
			}
		}
		return true
	}

	// Single values.
	item0 := arithmetic.ToArithmetic(value0.Content[0].Data)
	item1 := arithmetic.ToArithmetic(value1.Content[0].Data)
	switch operator {
	case grammar.Equals: // Equals.
		if item0 != item1 {
			return false
		}
	case grammar.NotEquals: // Not equals.
		if item0 == item1 {
			return false
		}
	case grammar.TokenGreat: // Greater.
		if item0 <= item1 {
			return false
		}
	case grammar.TokenLess: // Less.
		if item0 >= item1 {
			return false
		}
	case grammar.GreaterEquals: // Greater or equals.
		if item0 < item1 {
			return false
		}
	case grammar.LessEquals: // Less or equals.
		if item0 > item1 {
			return false
		}
	}

	return true
}

// processCondition Process conditional expression and return result.
// tokens Tokens to process.
func (i *Interpreter) processCondition(tokens *[]obj.Token) string {
	i.processRange(tokens)

	// Process condition.
	ors := parser.DecomposeConditionalProcess(*tokens, grammar.TokenVerticalBar)
	for _, current := range *ors {
		// Decompose and conditions.
		ands := parser.DecomposeConditionalProcess(current, grammar.TokenAmper)
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
