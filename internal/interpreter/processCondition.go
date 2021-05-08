package interpreter

import (
	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

func compareValues(operator string, data0, data1 objects.DataFrame) bool {
	if data0.Type != data1.Type &&
		(data0.Type == fract.VALString || data1.Type == fract.VALString) {
		return false
	}

	switch operator {
	case grammar.Equals: // Equals.
		if (data0.Type == fract.VALString && data0.Data != data1.Data) ||
			(data0.Type != fract.VALString &&
				arithmetic.ToArithmetic(data0.Data) != arithmetic.ToArithmetic(data1.Data)) {
			return false
		}
	case grammar.NotEquals: // Not equals.
		if (data0.Type == fract.VALString && data0.Data == data1.Data) ||
			(data0.Type != fract.VALString &&
				arithmetic.ToArithmetic(data0.Data) == arithmetic.ToArithmetic(data1.Data)) {
			return false
		}
	case grammar.TokenGreat: // Greater.
		if (data0.Type == fract.VALString && data0.Data <= data1.Data) ||
			(data0.Type != fract.VALString &&
				arithmetic.ToArithmetic(data0.Data) <= arithmetic.ToArithmetic(data1.Data)) {
			return false
		}
	case grammar.TokenLess: // Less.
		if (data0.Type == fract.VALString && data0.Data >= data1.Data) ||
			(data0.Type != fract.VALString &&
				arithmetic.ToArithmetic(data0.Data) >= arithmetic.ToArithmetic(data1.Data)) {
			return false
		}
	case grammar.GreaterEquals: // Greater or equals.
		if (data0.Type == fract.VALString && data0.Data < data1.Data) ||
			(data0.Type != fract.VALString &&
				arithmetic.ToArithmetic(data0.Data) < arithmetic.ToArithmetic(data1.Data)) {
			return false
		}
	case grammar.LessEquals: // Less or equals.
		if (data0.Type == fract.VALString && data0.Data > data1.Data) ||
			(data0.Type != fract.VALString &&
				arithmetic.ToArithmetic(data0.Data) > arithmetic.ToArithmetic(data1.Data)) {
			return false
		}
	}
		return true
}

func compare(value0, value1 objects.Value, operator string) bool {
	// String comparison.
	if !value0.Array || !value1.Array {
		data0 := value0.Content[0]
		data1 := value1.Content[0]
		if (data0.Type == fract.VALString && data1.Type != fract.VALString) ||
			(data0.Type != fract.VALString && data1.Type == fract.VALString) {
			return false
		}
		return compareValues(operator, data0, data1)
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

		for index, val0Content := range value0.Content {
			if !compareValues(operator, val0Content, value1.Content[index]) {
				return false
			}
		}

		return true
	}

	// Single value comparison.
	return compareValues(operator, value0.Content[0], value1.Content[0])
}

// processCondition returns condition result.
func (i *Interpreter) processCondition(tokens []objects.Token) string {
	i.processRange(&tokens)

	TRUE := objects.Value{Content: []objects.DataFrame{{Data: grammar.KwTrue}}}

	// Process condition.
	ors := parser.DecomposeConditionalProcess(tokens, grammar.LogicalOr)
	for _, or := range *ors {
		// Decompose and conditions.
		ands := parser.DecomposeConditionalProcess(or, grammar.LogicalAnd)
		// Is and long statement?
		if len(*ands) > 1 {
			for _, and := range *ands {
				operatorIndex, operator := parser.FindConditionOperator(and)

				// Operator is not found?
				if operatorIndex == -1 {
					if compare(i.processValue(and), TRUE, grammar.Equals) {
						return grammar.KwTrue
					}
					continue
				}

				// Operator is first or last?
				if operatorIndex == 0 {
					fract.Error(and[0], "Comparison values are missing!")
				} else if operatorIndex == len(and)-1 {
					fract.Error(and[len(and)-1], "Comparison values are missing!")
				}

				if !compare(
					i.processValue(*vector.Sublist(and, 0, operatorIndex)),
					i.processValue(*vector.Sublist(and, operatorIndex+1, len(and)-operatorIndex-1)),
					operator) {
					return grammar.KwFalse
				}
			}
			return grammar.KwTrue
		}

		operatorIndex, operator := parser.FindConditionOperator(or)

		// Operator is not found?
		if operatorIndex == -1 {
			if compare(i.processValue(or), TRUE, grammar.Equals) {
				return grammar.KwTrue
			}
			continue
		}

		// Operator is first or last?
		if operatorIndex == 0 {
			fract.Error(or[0], "Comparison values are missing!")
		} else if operatorIndex == len(or)-1 {
			fract.Error(or[len(or)-1], "Comparison values are missing!")
		}

		if compare(
			i.processValue(*vector.Sublist(or, 0, operatorIndex)),
			i.processValue(*vector.Sublist(or, operatorIndex+1, len(or)-operatorIndex-1)),
			operator) {
			return grammar.KwTrue
		}
	}

	return grammar.KwFalse
}
