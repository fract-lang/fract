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
				}
				return data.Data
			}
		}
	}
	return data.Data
}
