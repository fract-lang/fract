package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// procMacro process macros and returns keyword state.
func (p *Parser) procMacro(tks []obj.Token) uint8 {
	tks = tks[1:]
	switch tks[0].T {
	case fract.Name:
		switch tks[0].V {
		case "pragma":
			if len(tks) != 2 || tks[1].T != fract.Name {
				fract.IPanic(tks[0], obj.SyntaxPanic, "Invalid pragma syntax!")
			}
			switch tks[1].V {
			case "enofi":
				if p.loopCount == -1 {
					p.loopCount = 0
				}
			default:
				fract.IPanic(tks[1], obj.SyntaxPanic, "Invalid pragma!")
			}
		default:
			fract.IPanic(tks[0], obj.SyntaxPanic, "Invalid macro!")
		}
	default:
		fract.IPanic(tks[0], obj.SyntaxPanic, "Invalid macro!")
	}
	return fract.None
}
