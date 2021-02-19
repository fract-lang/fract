/*
	skipBlock Function.
*/

package interpreter

import (
	"../fract"
	"../objects"
	"../utilities/vector"
)

// skipBlock Skip to block end.
// tokens First tokens.
func (i *Interpreter) skipBlock(tokens *vector.Vector) {
	blockCount := 1
	for ; i.index < len(i.tokens.Vals); i.index++ {
		first := i.tokens.Vals[i.index].(*vector.Vector).Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd {
			blockCount--
			if blockCount == 0 {
				return
			}
		} else if first.Type == fract.TypeIf { // if-elif-else.
			blockCount++
		} else if first.Type == fract.TypeLoop { // Loop.
			blockCount++
		}
	}
}
