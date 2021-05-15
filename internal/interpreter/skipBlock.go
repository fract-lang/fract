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
		if first := tokens[0]; first.Type == fract.TypeBlockEnd {
			count--
		} else if first.Type == fract.TypeMacro {
			if parser.IsBlockStatement(tokens) {
				count++
			} else if tokens[1].Type == fract.TypeBlockEnd {
				count--
			}
		} else if parser.IsBlockStatement(tokens) {
			count++
		}

		if count == 0 { return }
	}
}
