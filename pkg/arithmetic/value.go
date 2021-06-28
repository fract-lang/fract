package arithmetic

import "strconv"

// Value parse to arithmetic value.
func Value(v string) float64 {
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
