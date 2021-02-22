/*
	SolveArithmeticProcess Function.
*/

package arithmetic

import (
	"fmt"
	"math"

	fract ".."
	"../../grammar"
	"../../objects"
)

// solve Solve process.
// operator Operator.
// first First value.
// second Second value.
func solve(operator objects.Token, first float64, second float64) float64 {
	var result float64

	if operator.Value == grammar.TokenBackslash ||
		operator.Value == grammar.IntegerDivideWithBigger { // Divide with bigger.
		if operator.Value == grammar.TokenBackslash {
			operator.Value = grammar.TokenSlash
		} else {
			operator.Value = grammar.IntegerDivision
		}

		if first < second {
			cache := first
			first = second
			second = cache
		}
	}

	if operator.Value == grammar.TokenPlus { // Addition.
		result = first + second
	} else if operator.Value == grammar.TokenMinus { // Subtraction.
		result = first - second
	} else if operator.Value == grammar.TokenStar { // Multiply.
		result = first * second
	} else if operator.Value == grammar.TokenSlash ||
		operator.Value == grammar.IntegerDivision { // Division.
		if first == 0 || second == 0 {
			fract.Error(operator, "Divide by zero!")
		}
		result = first / second

		if operator.Value == grammar.IntegerDivision {
			result = math.RoundToEven(result)
		}
	} else if operator.Value == grammar.TokenCaret { // Exponentiation.
		result = math.Pow(first, second)
	} else if operator.Value == grammar.TokenPercent { // Mod.
		result = math.Mod(first, second)
	} else {
		fract.Error(operator, "Operator is invalid!")
	}

	return result
}

// SolveArithmeticProcess Solve arithmetic process.
// process Process to solve.
func SolveArithmeticProcess(process objects.ArithmeticProcess) objects.Value {
	var value objects.Value
	value.Type = fract.VTInteger

	value.Charray = process.FirstV.Charray || process.SecondV.Charray

	/* Check type. */
	if process.FirstV.Type == fract.VTFloat || process.SecondV.Type == fract.VTFloat {
		value.Type = fract.VTFloat
	}

	if process.FirstV.Array && process.SecondV.Array {
		if len(process.FirstV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		} else if len(process.SecondV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		}
		if len(process.FirstV.Content) != len(process.SecondV.Content) &&
			(len(process.FirstV.Content) != 1 && len(process.SecondV.Content) != 1) {
			fract.Error(process.Second, "Array element count is not one or equals to first array!")
		}

		if len(process.FirstV.Content) == 1 {
			first, _ := ToFloat64(process.FirstV.Content[0])
			for index := range process.SecondV.Content {
				second, _ := ToFloat64(process.SecondV.Content[index])
				process.SecondV.Content[index] = fmt.Sprintf("%g",
					solve(process.Operator, first, second))
			}
			value.Content = process.SecondV.Content
		} else if len(process.SecondV.Content) == 1 {
			second, _ := ToFloat64(process.SecondV.Content[0])
			for index := range process.FirstV.Content {
				first, _ := ToFloat64(process.FirstV.Content[index])
				process.FirstV.Content[index] = fmt.Sprintf("%g",
					solve(process.Operator, first, second))
			}
			value.Content = process.FirstV.Content
		} else {
			for index := range process.FirstV.Content {
				first, _ := ToFloat64(process.FirstV.Content[index])
				second, _ := ToFloat64(process.SecondV.Content[index])
				process.FirstV.Content[index] = fmt.Sprintf("%g",
					solve(process.Operator, first, second))
			}
			value.Content = process.FirstV.Content
		}
		value.Array = true
	} else if process.FirstV.Array {
		if len(process.FirstV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		}

		second, _ := ToFloat64(process.SecondV.Content[0])
		for index := range process.FirstV.Content {
			first, _ := ToFloat64(process.FirstV.Content[index])
			process.FirstV.Content[index] = fmt.Sprintf("%g",
				solve(process.Operator, first, second))
		}
		value.Array = true
		value.Content = process.FirstV.Content
	} else if process.SecondV.Array {
		if len(process.SecondV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		}

		first, _ := ToFloat64(process.FirstV.Content[0])
		for index := range process.SecondV.Content {
			second, _ := ToFloat64(process.SecondV.Content[index])
			process.SecondV.Content[index] = fmt.Sprintf("%g",
				solve(process.Operator, second, first))
		}
		value.Array = true
		value.Content = process.SecondV.Content
	} else {
		first, _ := ToFloat64(process.FirstV.Content[0])
		second, _ := ToFloat64(process.SecondV.Content[0])
		value.Content = []string{fmt.Sprintf("%g",
			solve(process.Operator, first, second))}
	}

	return value
}
