package value

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

// Val instance.
type Val struct {
	D interface{}
	T uint8
}

// Returns immutable copy.
func (d Val) Immut() Val {
	v := Val{T: d.T}
	switch d.T {
	case Map:
		c := MapModel{}
		for k, v := range d.D.(MapModel) {
			c[k] = v
		}
		v.D = c
	case Array:
		c := make([]Val, len(d.D.([]Val)))
		copy(c, d.D.([]Val))
		v.D = c
	default:
		v.D = d.D
	}
	return v
}

func (d Val) String() string {
	switch d.T {
	case Func:
		return "object.func"
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
	case Str, Array, Map:
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
		return len(v.D.(MapModel))
	}
	return 0
}

func (v Val) Equals(dt Val) bool {
	return reflect.DeepEqual(v.D, dt.D)
}

func (v Val) NotEquals(dt Val) bool {
	return !v.Equals(dt)
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
