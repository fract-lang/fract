/*
	FormatData Function.
*/

package fract

import (
	"math/big"

	obj "github.com/fract-lang/fract/pkg/objects"
)

// FormatData Format data value.
// data Data to format.
func FormatData(data obj.DataFrame) string {
	if data.Type == VALString {
		return data.Data
	}

	for index := len(data.Data) - 1; index >= 0; index-- {
		if ch := data.Data[index]; ch != '0' {
			if ch == '.' {
				if data.Type == VALFloat {
					data.Data = data.Data[:index+1]
					data.Data += "0"
				} else {
					if index == 0 {
						index++
					}
					data.Data = data.Data[:index]
				}
			} else if data.Type == VALFloat {
				data.Data = data.Data[:index+1]
			}
			break
		}
	}

	if data.Type == VALFloat {
		b := new(big.Float)
		b.SetString(data.Data)
		data.Data = b.String()
	}

	return data.Data
}
