package fract

import (
	"math/big"
	"strings"

	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

func FormatData(data objects.DataFrame) string {
	if data.Type == VALString || data.Type == VALBoolean {
		return data.Data
	}

	if data.Data != grammar.KwNaN {
		if data.Type == VALInteger {
			bigfloat, _ := new(big.Float).SetString(data.Data)
			data.Data = bigfloat.String()
			return data.Data
		}

		b, _ := new(big.Float).SetString(data.Data)
		data.Data = b.String()
		if !strings.Contains(data.Data, ".") {
			data.Data += ".0"
		}
	}

	return data.Data
}
