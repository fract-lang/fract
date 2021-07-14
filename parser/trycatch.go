package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// procTryCatch process try-catch blocks and returns keyword state.
func (p *Parser) procTryCatch(tks obj.Tokens) uint8 {
	fract.TryCount++
	var (
		vlen = len(p.vars)
		flen = len(p.funcs)
		ilen = len(p.Imports)
		dlen = len(defers)
		kws  = fract.None
	)
	(&obj.Block{
		Try: func() {
			for _, tks := range p.getBlock(tks[1:]) {
				if kws = p.process(tks); kws != fract.None {
					break
				}
			}
			if p.Tks[p.i+1][0].T == fract.Catch {
				p.i++
			}
			fract.TryCount--
			p.vars = p.vars[:vlen]
			p.funcs = p.funcs[:flen]
			p.Imports = p.Imports[:ilen]
			for index := len(defers) - 1; index >= dlen; index-- {
				defers[index].call()
			}
			defers = defers[:dlen]
		},
		Catch: func(cp obj.Panic) {
			p.loopCount = 0
			fract.TryCount--
			p.vars = p.vars[:vlen]
			p.funcs = p.funcs[:flen]
			p.Imports = p.Imports[:ilen]
			defers = defers[:dlen]
			p.i++
			tks = p.Tks[p.i]
			if tks[0].T != fract.Catch {
				p.i--
				return
			}
			for _, tks := range p.getBlock(tks[1:]) {
				if kws = p.process(tks); kws != fract.None {
					break
				}
			}
			p.vars = p.vars[:vlen]
			p.funcs = p.funcs[:flen]
			for i := len(defers) - 1; i >= dlen; i-- {
				defers[i].call()
			}
			defers = defers[:dlen]
		},
	}).Do()
	return kws
}
