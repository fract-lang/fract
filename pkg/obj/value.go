package obj

import (
	"fmt"
	"math/big"
	"strings"
)

const (
	VInt   uint8 = 0
	VFloat uint8 = 1
	VStr   uint8 = 2
	VBool  uint8 = 3
	VFunc  uint8 = 4
	VArray uint8 = 5
)

func stringArray(src []Data) string {
	if len(src) == 0 {
		return "[]"
	}
	var sb strings.Builder
	sb.WriteByte('[')
	for _, d := range src {
		sb.WriteString(d.Format() + " ")
	}
	return sb.String()[:sb.Len()-1] + "]"
}

// Data instance.
type Data struct {
	D interface{}
	T uint8
}

// Get data as string.
func (d Data) String() string {
	switch d.T {
	case VFunc:
		return "object.function"
	case VArray:
		return stringArray(d.D.([]Data))
	default:
		if d.D == nil {
			return ""
		}
		return d.D.(string)
	}
}

func (d Data) Format() string {
	if d.T != VInt && d.T != VFloat {
		return d.String()
	}
	dt := d.String()
	if dt != "NaN" {
		if d.T == VInt {
			bf, _ := new(big.Float).SetString(dt)
			dt = bf.String()
			return dt
		}
		b, _ := new(big.Float).SetString(dt)
		dt = b.String()
		if !strings.Contains(dt, ".") {
			dt = dt + ".0"
		}
	}
	return dt
}

// Value intance.
type Value struct {
	D   []Data
	Arr bool
}

func (v *Value) String() string {
	if v.D == nil {
		return ""
	}
	if v.Arr {
		return stringArray(v.D)
	}
	return v.D[0].Format()
}

func (v *Value) Print() bool {
	if v.D == nil {
		return false
	}
	fmt.Print(v.String())
	return true
}
