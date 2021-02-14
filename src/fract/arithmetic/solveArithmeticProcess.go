/*
	SolveArithmeticProcess Function.
*/

package arithmetic

import (
	"math"

	fract ".."
	"../../grammar"
	"../../objects"
)

// SolveArithmeticProcess Solve arithmetic process.
// process Process to solve.
func SolveArithmeticProcess(process objects.ArithmeticProcess) (int, float64) {
	/* Check type. */
	_type := fract.VTInteger
	if IsFloatValue(process.First.Value) || IsFloatValue(process.Second.Value) {
		_type = fract.VTFloat

		if IsFloatValue(process.First.Value) && !CheckFloat(process.First.Value) {
			fract.Error(process.First, "Decimal limit exceeded!")
		}
		if IsFloatValue(process.Second.Value) && !CheckFloat(process.Second.Value) {
			fract.Error(process.Second, "Decimal limit exceeded!")
		}
	}

	var result float64

	first, err := ToFloat64(process.First.Value)
	if err != nil {
		fract.Error(process.First, "Value out of range!")
	}
	second, err := ToFloat64(process.Second.Value)
	if err != nil {
		fract.Error(process.Second, "Value out of range!")
	}

	if process.Operator.Value == grammar.TokenReverseSlash ||
		process.Operator.Value == grammar.IntegerDivideWithBigger { // Divide with bigger.
		if process.Operator.Value == grammar.TokenReverseSlash {
			process.Operator.Value = grammar.TokenSlash
		} else {
			process.Operator.Value = grammar.IntegerDivision
		}

		if first < second {
			cache := first
			first = second
			second = cache
		}
	}

	if process.Operator.Value == grammar.TokenPlus { // Addition.
		result = first + second
	} else if process.Operator.Value == grammar.TokenMinus { // Subtraction.
		result = first - second
	} else if process.Operator.Value == grammar.TokenStar { // Multiply.
		result = first * second
	} else if process.Operator.Value == grammar.TokenSlash ||
		process.Operator.Value == grammar.IntegerDivision { // Division.
		if first == 0 {
			fract.Error(process.First, "Divide by zero!")
		} else if second == 0 {
			fract.Error(process.Second, "Divide by zero!")
		}
		result = first / second

		if process.Operator.Value == grammar.IntegerDivision {
			result = math.RoundToEven(result)
		}
	} else if process.Operator.Value == grammar.TokenCaret { // Exponentiation.
		result = math.Pow(first, second)
	} else if process.Operator.Value == grammar.TokenPercent { // Mod.
		result = math.Mod(first, second)
	} else {
		fract.Error(process.Operator,
			"Operator is invalid!: "+process.Operator.Value)
	}

	return _type, result
}
