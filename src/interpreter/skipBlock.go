/*
	skipBlock Function.
*/

package interpreter

import (
	"../fract"
	"../objects"
	"../utils/vector"
)

// skipBlock Skip to block end.
func (i *Interpreter) skipBlock() {
	count := 1
	for ; i.index < len(i.tokens.Vals); i.index++ {
		first := i.tokens.Vals[i.index].(vector.Vector).Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd {
			count--
			if count == 0 {
				return
			}
		} else if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
			first.Type == fract.TypeFunction {
			count++
		}
	}
}
