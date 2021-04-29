/*
	FormatData Function.
*/

package fract

import (
	"math/big"
	"strings"

	obj "github.com/fract-lang/fract/pkg/objects"
)

// FormatData Format data value.
// data Data to format.
func FormatData(data obj.DataFrame) string {
	if data.Type == VALString || data.Type == VALBoolean {
		return data.Data
	}

	if data.Data != "NaN" {
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
