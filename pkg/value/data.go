package value

import (
	"math/big"
	"strings"
)

// Data instance.
type Data struct {
	D interface{}
	T uint8
}

func (d Data) String() string {
	switch d.T {
	case Func:
		return "object.function"
	case Array:
		return stringArray(d.D.([]Data))
	default:
		if d.D == nil {
			return ""
		}
		return d.D.(string)
	}
}

func (d Data) Format() string {
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

func (d Data) Equals(dt Data) bool {
	return (d.T == Str && d.D == dt.D) || (d.T != Str && Conv(d.String()) == Conv(dt.String()))
}

func (d Data) NotEquals(dt Data) bool {
	return (d.T == Str && d.D != dt.D) || (d.T != Str && Conv(d.String()) != Conv(dt.String()))
}

func (d Data) Greater(dt Data) bool {
	return (d.T == Str && d.String() > dt.String()) || (d.T != Str && Conv(d.String()) > Conv(dt.String()))
}

func (d Data) Less(dt Data) bool {
	return (d.T == Str && d.String() < dt.String()) || (d.T != Str && Conv(d.String()) < Conv(dt.String()))
}

func (d Data) GreaterEquals(dt Data) bool {
	return (d.T == Str && d.String() >= dt.String()) || (d.T != Str && Conv(d.String()) >= Conv(dt.String()))
}

func (d Data) LessEquals(dt Data) bool {
	return (d.T == Str && d.String() <= dt.String()) || (d.T != Str && Conv(d.String()) <= Conv(dt.String()))
}
