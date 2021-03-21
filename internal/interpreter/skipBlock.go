/*
	skipBlock Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// skipBlock Skip to block end.
// ifBlock Enable skip if statement is block start?
func (i *Interpreter) skipBlock(ifBlock bool) {
	if ifBlock {
		first := i.tokens.Vals[i.index].(vector.Vector).Vals[0].(obj.Token)
		if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
			first.Type == fract.TypeFunction {
			i.index++
		} else {
			return
		}
	}

	count := 1
	for ; i.index < len(i.tokens.Vals); i.index++ {
		tokens := i.tokens.Vals[i.index].(vector.Vector)
		if first := tokens.Vals[0].(obj.Token); first.Type == fract.TypeBlockEnd {
			count--
			if count == 0 {
				return
			}
		} else if parser.IsBlockStatement(tokens) {
			count++
		}
	}
}
