package arithmetic

import (
	"strconv"
)

// ToArithmetic parse value to arithmetic value.
func ToArithmetic(value string) float64 {
	switch value {
	case "true":
		return 1
	case "false":
		return 0
	default:
		flt, _ := strconv.ParseFloat(value, 64)
		return flt
	}
}
