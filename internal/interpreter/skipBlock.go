/*
	skipBlock Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/parser"
)

// skipBlock Skip to block end.
// ifBlock Enable skip if statement is block start?
func (i *Interpreter) skipBlock(ifBlock bool) {
	if ifBlock {
		first := i.tokens[i.index][0]
		if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
			first.Type == fract.TypeFunction {
			i.index++
		} else {
			return
		}
	}

	count := 1
	for ; i.index < len(i.tokens); i.index++ {
		tokens := i.tokens[i.index]
		if tokens[0].Type == fract.TypeBlockEnd {
			count--
			if count == 0 {
				return
			}
		} else if parser.IsBlockStatement(tokens) {
			count++
		}
	}
}
