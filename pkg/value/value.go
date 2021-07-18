package value

import (
	"strconv"
)

const (
	Int   uint8 = 1
	Float uint8 = 2
	Str   uint8 = 3
	Bool  uint8 = 4
	Func  uint8 = 5
	Array uint8 = 6
	Map   uint8 = 7
)

type MapModel map[Val]Val

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
