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
func (i *Interpreter) skipBlock() {
	blockCount := 1
	for ; i.index < len(i.tokens.Vals); i.index++ {
		first := i.tokens.Vals[i.index].(vector.Vector).Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd {
			blockCount--
			if blockCount == 0 {
				return
			}
		} else if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
			first.Type == fract.TypeFunction {
			blockCount++
		}
	}

	if blockCount > 0 { // Check blocks.
		i.lexer.Line--
		i.lexer.Error("Block is expected ending...")
	}
}
