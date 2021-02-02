package arithmetic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	fract ".."
	"../../grammar"
	"../../objects"
)

// IsNegative Is negative number?
// value Value to check.
func IsNegative(value string) bool {
	return strings.HasPrefix(value, grammar.TokenMinus)
}

// IsNumeric Char is numeric?
// char Char to check.
func IsNumeric(char byte) bool {
	return char == '0' ||
		char == '1' ||
		char == '2' ||
		char == '3' ||
		char == '4' ||
		char == '5' ||
		char == '6' ||
		char == '7' ||
		char == '8' ||
		char == '9'
}

// IsInteger Value is an integer?
// value Value to check.
func IsInteger(value string) bool {
	state, _ := regexp.MatchString("^(-|)\\s*[0-9]+$", value)
	return state
}

// IsFloat Value is an float?
// value Value to check.
func IsFloat(value string) bool {
	state, _ := regexp.MatchString("^(-|)\\s*[0-9]+(\\.[0-9]+)?$", value)
	return state
}

// ToFloat String to float.
// value Value to parse.
func ToFloat(value string) (float32, error) {
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(result), err
}

// ToDouble String to double.
// value Value to parse.
func ToDouble(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

// ToSByte String to 8bit integer.
// value Value to parse.
func ToSByte(value string) (int8, error) {
	result, err := strconv.ParseInt(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(result), err
}

// ToShort String to 16bit integer.
// value Value to parse.
func ToShort(value string) (int16, error) {
	result, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(result), err
}

// ToInt String to 32bit integer.
// value Value to parse.
func ToInt(value string) (int32, error) {
	result, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(result), err
}

// ToLong String to 64bit integer.
// value Value to parse.
func ToLong(value string) (int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return int64(result), err
}

// ToByte String to 8bit unsigned integer.
// value Value to parse.
func ToByte(value string) (uint8, error) {
	result, err := strconv.ParseUint(value, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(result), err
}

// ToUShort String to 16bit unsigned integer.
// value Value to parse.
func ToUShort(value string) (uint16, error) {
	result, err := strconv.ParseUint(value, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(result), err
}

// ToUInt String to 32bit integer.
// value Value to parse.
func ToUInt(value string) (uint32, error) {
	result, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(result), err
}

// ToULong String to 64bit unsigned integer.
// value Value to parse.
func ToULong(value string) (uint64, error) {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(result), err
}

// FloatToString Float to string.
// value Value to parse.
func FloatToString(value interface{}) string {
	return fmt.Sprintf("%f", value)
}

// IntToString Integer to string.
// value Value to parse.
func IntToString(value interface{}) string {
	return fmt.Sprintf("%d", value)
}

// SolveArithmeticProcess Solve arithmetic process.
// process Process to solve.
func SolveArithmeticProcess(process objects.ArithmeticProcess) float64 {
	var result float64

	first, _ := ToDouble(process.First.Value)
	second, _ := ToDouble(process.Second.Value)

	if process.Operator.Value == grammar.TokenPlus {
		result = first + second
	} else if process.Operator.Value == grammar.TokenMinus {
		result = first - second
	} else if process.Operator.Value == grammar.TokenStar {
		result = first * second
	} else if process.Operator.Value == grammar.TokenSlash {
		if first == 0 {
			fract.Error(process.First, "Divide by zero!")
		} else if second == 0 {
			fract.Error(process.Second, "Divide by zero!")
		}
		result = first / second
	} else {
		fract.Error(process.Operator,
			"Operator is invalid!: "+process.Operator.Value)
	}

	return result
}
