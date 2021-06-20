package obj

import (
	"fmt"
	"math/big"
	"strings"
)

// TODO: Minimize here.
// TODO: Add []Datas to array string function.

const (
	VInteger  uint8 = 0
	VFloat    uint8 = 1
	VString   uint8 = 2
	VBoolean  uint8 = 3
	VFunction uint8 = 4
	VArray    uint8 = 5
)

// Data instance.
type Data struct {
	D interface{}
	T uint8
}

// Get data as string.
func (d Data) String() string {
	switch d.T {
	case VFunction:
		return "object.function"
	case VArray:
		if len(d.D.([]Data)) == 0 {
			return "[]"
		} else {
			var sb strings.Builder
			sb.WriteByte('[')
			for _, data := range d.D.([]Data) {
				sb.WriteString(data.Format() + " ")
			}
			return sb.String()[:sb.Len()-1] + "]"
		}
	default:
		if d.D == nil {
			return "0"
		}
		return d.D.(string)
	}
}

func (d Data) Format() string {
	if d.T == VString || d.T == VBoolean || d.T == VFunction || d.T == VArray {
		return d.String()
	}
	dt := d.String()
	if dt != "NaN" {
		if d.T == VInteger {
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

func (v Value) String() string {
	if v.D == nil {
		return ""
	}

	if v.Arr {
		if len(v.D) == 0 {
			return "[]"
		}
		var sb strings.Builder
		sb.WriteByte('[')
		for _, data := range v.D {
			sb.WriteString(data.Format() + " ")
		}
		return sb.String()[:sb.Len()-1] + "]"
	}
	return v.D[0].Format()
}

func (v *Value) Print() bool {
	if v.D == nil {
		return false
	}

	if v.Arr {
		if len(v.D) == 0 {
			fmt.Print("[]")
		} else {
			var sb strings.Builder
			sb.WriteByte('[')
			for _, data := range v.D {
				sb.WriteString(data.Format() + " ")
			}
			fmt.Print(sb.String()[:sb.Len()-1] + "]")
		}
	} else {
		fmt.Print(v.D[0].Format())
	}
	return true
}
