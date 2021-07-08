package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// procIf process if-elif-else returns keyword state.
func (p *Parser) procIf(tks obj.Tokens) uint8 {
	tkslen := len(tks)
	ctks := tks.Sub(1, tkslen-1)
	// Condition is empty?
	if ctks == nil {
		first := tks[0]
		fract.Errorc(first.F, first.Ln, first.Col+len(first.Val), "Condition is empty!")
	}
	s := p.procCondition(*ctks)
	vlen := len(p.vars)
	flen := len(p.funcs)
	kws := fract.None
	/* Interpret/skip block. */
	for {
		p.i++
		tks := p.Tks[p.i]
		first := tks[0]
		if first.T == fract.End { // Block is ended.
			goto end
		} else if first.T == fract.ElseIf { // Else if block.
			tkslen = len(tks)
			ctks := tks.Sub(1, tkslen-1)
			// Condition is empty?
			if ctks == nil {
				first := tks[0]
				fract.Errorc(first.F, first.Ln, first.Col+len(first.Val), "Condition is empty!")
			}
			if s == "true" {
				p.skipBlock(false)
				goto end
			}
			s = p.procCondition(*ctks)
			// Interpret/skip block.
			for {
				p.i++
				tks := p.Tks[p.i]
				first := tks[0]
				if first.T == fract.End { // Block is ended.
					goto end
				} else if first.T == fract.If { // If block.
					if s == "true" && kws == fract.None {
						p.procIf(tks)
					} else {
						p.skipBlock(true)
					}
					continue
				} else if first.T == fract.ElseIf || first.T == fract.Else { // Else if or else block.
					p.i--
					break
				}
				// Condition is true?
				if s == "true" && kws == fract.None {
					if kws = p.process(tks); kws != fract.None {
						p.skipBlock(false)
					}
				} else {
					p.skipBlock(true)
				}
			}
			if s == "true" {
				p.skipBlock(false)
				goto end
			}
			continue
		} else if first.T == fract.Else { // Else block.
			if len(tks) > 1 {
				fract.Error(first, "Else block is not take any arguments!")
			}
			if s == "true" {
				p.skipBlock(false)
				goto end
			}
			/* Interpret/skip block. */
			for {
				p.i++
				tks := p.Tks[p.i]
				first := tks[0]
				if first.T == fract.End { // Block is ended.
					goto end
				} else if first.T == fract.If { // If block.
					if kws == fract.None {
						p.procIf(tks)
					} else {
						p.skipBlock(true)
					}
					continue
				}
				// Condition is true?
				if kws == fract.None {
					if kws = p.process(tks); kws != fract.None {
						p.skipBlock(false)
					}
				}
			}
		}
		// Condition is true?
		if s == "true" && kws == fract.None {
			if kws = p.process(tks); kws != fract.None {
				p.skipBlock(false)
			}
		} else {
			p.skipBlock(true)
		}
	}
end:
	p.vars = p.vars[:vlen]
	p.funcs = p.funcs[:flen]
	return kws
}
