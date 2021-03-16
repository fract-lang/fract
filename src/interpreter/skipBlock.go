/*
	skipBlock Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// skipBlock Skip to block end.
// ifBlock Enable skip if statement is block start?
func (i *Interpreter) skipBlock(ifBlock bool) {
	if ifBlock {
		first := i.tokens.Vals[i.index].(vector.Vector).Vals[0].(objects.Token)
		if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
			first.Type == fract.TypeFunction {
			i.index++
		} else {
			return
		}
	}

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
