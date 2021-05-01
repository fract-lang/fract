package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/parser"
)

// skipBlock skip to block end.
func (i *Interpreter) skipBlock(ifBlock bool) {
	if ifBlock {
		if parser.IsBlockStatement(i.Tokens[i.index]) {
			i.index++
		} else {
			return
		}
	}

	count := 1
	i.index--
	for {
		i.index++
		tokens := i.Tokens[i.index]
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
