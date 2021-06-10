package objects

import (
	"math/big"
	"strings"

	"github.com/fract-lang/fract/pkg/grammar"
)

const (
	VALInteger  uint8 = 0
	VALFloat    uint8 = 1
	VALString   uint8 = 2
	VALBoolean  uint8 = 3
	VALFunction uint8 = 4
)

// Data instance.
type Data struct {
	Data interface{}
	Type uint8
}

// Get data as string.
func (d Data) String() string {
	switch d.Type {
	case VALFunction:
		return "object.function"
	default:
		return d.Data.(string)
	}
}

func (d Data) Format() string {
	data := d.String()
	if d.Type == VALString || d.Type == VALBoolean || d.Type == VALFunction {
		return data
	}
	if data != grammar.KwNaN {
		if d.Type == VALInteger {
			bigfloat, _ := new(big.Float).SetString(data)
			data = bigfloat.String()
			return data
		}
		b, _ := new(big.Float).SetString(data)
		data = b.String()
		if !strings.Contains(data, ".") {
			data = data + ".0"
		}
	}
	return data
}
