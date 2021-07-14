package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// procIf process if-elif-else returns keyword state.
func (p *Parser) procIf(tks obj.Tokens) uint8 {
	bi := findBlock(tks)
	btks, ctks := p.getBlock(tks[bi:]), tks[1:bi]
	// Condition is empty?
	if len(ctks) == 0 {
		first := tks[0]
		fract.IPanicC(first.F, first.Ln, first.Col+len(first.V), obj.SyntaxPanic, "Condition is empty!")
	}
	s := p.procCondition(ctks)
	vlen := len(p.vars)
	flen := len(p.funcs)
	ilen := len(p.Imports)
	kws := fract.None
	for _, tks := range btks {
		// Condition is true?
		if s == "true" && kws == fract.None {
			if kws = p.process(tks); kws != fract.None {
				break
			}
		} else {
			break
		}
	}
rep:
	p.i++
	if p.i >= len(p.Tks) {
		p.i--
		goto end
	}
	tks = p.Tks[p.i]
	if tks[0].T != fract.Else {
		p.i--
		goto end
	}
	if len(tks) > 1 && tks[1].T == fract.If { // Else if.
		bi = findBlock(tks)
		btks, ctks = p.getBlock(tks[bi:]), tks[2:bi]
		// Condition is empty?
		if len(ctks) == 0 {
			first := tks[1]
			fract.IPanicC(first.F, first.Ln, first.Col+len(first.V), obj.SyntaxPanic, "Condition is empty!")
		}
		if s == "true" {
			goto rep
		}
		s = p.procCondition(ctks)
		for _, tks := range btks {
			// Condition is true?
			if s == "true" && kws == fract.None {
				if kws = p.process(tks); kws != fract.None {
					break
				}
			} else {
				break
			}
		}
		goto rep
	}
	btks = p.getBlock(tks[1:])
	if s == "true" {
		goto end
	}
	for _, tks := range btks {
		// Condition is true?
		if kws == fract.None {
			if kws = p.process(tks); kws != fract.None {
				break
			}
		}
	}
end:
	p.vars = p.vars[:vlen]
	p.funcs = p.funcs[:flen]
	p.Imports = p.Imports[:ilen]
	return kws
}
