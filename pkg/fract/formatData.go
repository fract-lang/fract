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
	if data.Type == VALString {
		return data.Data
	}

	b := new(big.Float)
	if _, fail := b.SetString(data.Data); !fail {
		data.Data = b.String()
		if data.Type == VALFloat && !strings.Contains(data.Data, ".") {
			data.Data += ".0"
		}
	}

	return data.Data
}
