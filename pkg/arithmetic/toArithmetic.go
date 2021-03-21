package arithmetic

import (
	"github.com/fract-lang/fract/pkg/grammar"
)

// ToArithmetic Parse value to arithmetic value.
// value Value to parse.
func ToArithmetic(value string) float64 {
	if value == grammar.KwTrue {
		return 1
	} else if value == grammar.KwFalse {
		return 0
	}
	flt, _ := ToFloat64(value)
	return flt
}
