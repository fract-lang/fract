/*
	FormatData Function.
*/

package fract

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// FormatData Format data value.
// data Data to format.
func FormatData(data obj.DataFrame) string {
	if data.Type != VALString {
		for index := len(data.Data) - 1; index >= 0; index-- {
		repeat:
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

				if data.Type != VALFloat {
					length := len(data.Data)
					if length > 2 && data.Data[length-2:] != ".0" {
						data.Type = VALFloat
						goto repeat
					}
				}
				return data.Data
			}
		}
	}
	return data.Data
}
