package value

import (
	"fmt"
	"math/big"
	"strings"
)

// Val instance.
type Val struct {
	D interface{}
	T uint8
}

func (d Val) String() string {
	switch d.T {
	case Func:
		return "object.function"
	case Array:
		return fmt.Sprint(d.D)
	case Map:
		s := fmt.Sprint(d.D)
		return "{" + s[4:len(s)-1] + "}"
	default:
		if d.D == nil {
			return ""
		}
		return d.D.(string)
	}
}

func (d Val) Format() string {
	if d.T != Int && d.T != Float {
		return d.String()
	}
	dt := d.String()
	if dt != "NaN" {
		if d.T == Int {
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

func (v Val) Print() bool {
	if v.D == nil {
		return false
	}
	fmt.Print(v.String())
	return true
}

// Is enumerable?
func (v Val) IsEnum() bool {
	switch v.T {
	case Array, Map:
		return true
	default:
		return false
	}
}

// Length.
func (v Val) Len() int {
	switch v.T {
	case Str:
		return len(v.D.(string))
	case Array:
		return len(v.D.([]Val))
	case Map:
		return len(v.D.(map[interface{}]Val))
	}
	return 0
}

func (v Val) Equals(dt Val) bool {
	return (v.T == Str && v.D == dt.D) || (v.T != Str && Conv(v.String()) == Conv(dt.String()))
}

func (v Val) NotEquals(dt Val) bool {
	return (v.T == Str && v.D != dt.D) || (v.T != Str && Conv(v.String()) != Conv(dt.String()))
}

func (v Val) Greater(dt Val) bool {
	return (v.T == Str && v.String() > dt.String()) || (v.T != Str && Conv(v.String()) > Conv(dt.String()))
}

func (v Val) Less(dt Val) bool {
	return (v.T == Str && v.String() < dt.String()) || (v.T != Str && Conv(v.String()) < Conv(dt.String()))
}

func (v Val) GreaterEquals(dt Val) bool {
	return (v.T == Str && v.String() >= dt.String()) || (v.T != Str && Conv(v.String()) >= Conv(dt.String()))
}

func (v Val) LessEquals(dt Val) bool {
	return (v.T == Str && v.String() <= dt.String()) || (v.T != Str && Conv(v.String()) <= Conv(dt.String()))
}
