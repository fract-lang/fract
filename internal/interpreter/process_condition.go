package interpreter

import (
	"strings"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

func compareValues(operator string, data0, data1 objects.Data) bool {
	if data0.Type != data1.Type && (data0.Type == objects.VALString || data1.Type == objects.VALString) {
		return false
	}

	switch operator {
	case grammar.Equals: // Equals.
		if (data0.Type == objects.VALString && data0.Data != data1.Data) ||
			(data0.Type != objects.VALString && arithmetic.ToArithmetic(data0.String()) != arithmetic.ToArithmetic(data1.String())) {
			return false
		}
	case grammar.NotEquals: // Not equals.
		if (data0.Type == objects.VALString && data0.Data == data1.Data) ||
			(data0.Type != objects.VALString && arithmetic.ToArithmetic(data0.String()) == arithmetic.ToArithmetic(data1.String())) {
			return false
		}
	case ">": // Greater.
		if (data0.Type == objects.VALString && data0.String() <= data1.String()) ||
			(data0.Type != objects.VALString && arithmetic.ToArithmetic(data0.String()) <= arithmetic.ToArithmetic(data1.String())) {
			return false
		}
	case "<": // Less.
		if (data0.Type == objects.VALString && data0.String() >= data1.String()) ||
			(data0.Type != objects.VALString && arithmetic.ToArithmetic(data0.String()) >= arithmetic.ToArithmetic(data1.String())) {
			return false
		}
	case grammar.GreaterEquals: // Greater or equals.
		if (data0.Type == objects.VALString && data0.String() < data1.String()) ||
			(data0.Type != objects.VALString && arithmetic.ToArithmetic(data0.String()) < arithmetic.ToArithmetic(data1.String())) {
			return false
		}
	case grammar.LessEquals: // Less or equals.
		if (data0.Type == objects.VALString && data0.String() > data1.String()) ||
			(data0.Type != objects.VALString && arithmetic.ToArithmetic(data0.String()) > arithmetic.ToArithmetic(data1.String())) {
			return false
		}
	}
	return true
}

func compare(value0, value1 objects.Value, operator objects.Token) bool {
	// In.
	if operator.Value == grammar.KwIn {
		if !value1.Array && value1.Content[0].Type != objects.VALString {
			fract.Error(operator, "Value is not enumerable!")
		}
		if value1.Array {
			if value0.Array {
				for _, d := range value1.Content {
					for _, cd := range value0.Content {
						if compareValues(grammar.Equals, d, cd) {
							return true
						}
					}
				}
			} else {
				data := value0.Content[0].String()
				for _, d := range value1.Content {
					if strings.Contains(data, d.String()) {
						return true
					}
				}
			}
		} else {
			if value0.Array {
				data := value1.Content[0].String()
				for _, d := range value0.Content {
					if d.Type != objects.VALString {
						fract.Error(operator, "All datas is not string!")
					}
					if strings.Contains(data, d.String()) {
						return true
					}
				}
			} else {
				if value1.Content[0].Type != objects.VALString {
					fract.Error(operator, "All datas is not string!")
				}
				if strings.Contains(value1.Content[0].String(), value1.Content[0].String()) {
					return true
				}
			}
		}
		return false
	}
	// String comparison.
	if !value0.Array || !value1.Array {
		data0 := value0.Content[0]
		data1 := value1.Content[0]
		if (data0.Type == objects.VALString && data1.Type != objects.VALString) ||
			(data0.Type != objects.VALString && data1.Type == objects.VALString) {
			fract.Error(operator, "The in keyword should use with string or enumerable data types!")
		}
		return compareValues(operator.Value, data0, data1)
	}
	// Array comparison.
	if value0.Array || value1.Array {
		if (value0.Array && !value1.Array) || (!value0.Array && value1.Array) {
			return false
		}
		if len(value0.Content) != len(value1.Content) {
			return operator.Value == grammar.NotEquals
		}
		for index, val0Content := range value0.Content {
			if !compareValues(operator.Value, val0Content, value1.Content[index]) {
				return false
			}
		}
		return true
	}
	// Single value comparison.
	return compareValues(operator.Value, value0.Content[0], value1.Content[0])
}

// processCondition returns condition result.
func (i *Interpreter) processCondition(tokens []objects.Token) string {
	i.processRange(&tokens)
	TRUE := objects.Value{Content: []objects.Data{{Data: grammar.KwTrue}}}
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
					operator.Value = grammar.Equals
					if compare(i.processValue(and), TRUE, operator) {
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
			operator.Value = grammar.Equals
			if compare(i.processValue(or), TRUE, operator) {
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
