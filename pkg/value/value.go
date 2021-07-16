package value

import (
	"strconv"
)

const (
	Int   uint8 = 0
	Float uint8 = 1
	Str   uint8 = 2
	Bool  uint8 = 3
	Func  uint8 = 4
	Array uint8 = 5
	Map   uint8 = 6
)

// Parse string to arithmetic value.
func Conv(v string) float64 {
	switch v {
	case "true":
		return 1
	case "false":
		return 0
	default:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	}
}
