package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

func (p *Parser) procDel(tks []obj.Token) {
	tkslen := len(tks)
	// Value is not defined?
	if tkslen < 2 {
		first := tks[0]
		fract.Errorc(first.F, first.Ln, first.Col+len(first.Val), "Value is not defined!")
	}
	comma := false
	for j := 1; j < tkslen; j++ {
		t := tks[j]
		if comma {
			if t.T != fract.Comma {
				fract.Error(t, "Comma is not found!")
			}
			comma = false
			continue
		}
		// Token is not a deletable object?
		if t.T != fract.Name {
			fract.Error(t, "This is not deletable object!")
		}
		pos, src := p.varIndexByName(t)
		// Name is not defined?
		if pos == -1 {
			pos, src := p.funcIndexByName(t)
			if pos == -1 {
				fract.Error(t, "\""+t.Val+"\" is not defined!")
			}
			// Protected?
			if src.funcs[pos].Protected {
				fract.Error(t, "Protected objects cannot be deleted manually from memory!")
			}
			src.funcs = append(src.funcs[:pos], src.funcs[pos+1:]...)
			continue
		}
		// Protected?
		if src.vars[pos].Protected {
			fract.Error(t, "Protected objects cannot be deleted manually from memory!")
		}
		src.vars = append(src.vars[:pos], src.vars[pos+1:]...)
		comma = true
	}
}
