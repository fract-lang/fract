/*
	FLOAT FUNCTIONS
*/

package arithmetic

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"../../grammar"
)

// IsFloatValue Value is float?
// value Value to check.
func IsFloatValue(value string) bool {
	return strings.Index(value, grammar.TokenDot) != -1
}

// CheckFloat Check float value validate.
// value Value to check.
func CheckFloat(value string) bool {
	return len(value[strings.Index(value, grammar.TokenDot)+1:]) <= 6
}

// FloatToString Float to string.
// value Value to parse.
func FloatToString(value interface{}) string {
	return fmt.Sprintf("%f", value)
}

// IsFloat Value is an float?
// value Value to check.
func IsFloat(value string) bool {
	state, _ := regexp.MatchString("^(-|)\\s*[0-9]+(\\.[0-9]+)?$", value)
	return state
}

// ToFloat32 String to float.
// value Value to parse.
func ToFloat32(value string) (float32, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	result, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(result), err
}

// ToFloat64 String to double.
// value Value to parse.
func ToFloat64(value string) (float64, error) {
	if value == grammar.KwTrue {
		return 1, nil
	} else if value == grammar.KwFalse {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}
