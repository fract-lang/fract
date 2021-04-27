package arithmetic

import (
	"strconv"

	"github.com/fract-lang/fract/pkg/grammar"
)

// ToArithmetic Parse value to arithmetic value.
// value Value to parse.
func ToArithmetic(value string) float64 {
	switch value {
	case grammar.KwTrue:
		return 1
	case grammar.KwFalse:
		return 0
	default:
		flt, _ := strconv.ParseFloat(value, 64)
		return flt
	}
}
