package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// procTryCatch process try-catch blocks and returns keyword state.
func (p *Parser) procTryCatch(tks obj.Tokens) uint8 {
	if len(tks) > 1 {
		fract.Error(tks[1], "Invalid syntax!")
	}
	fract.TryCount++
	var (
		vlen = len(p.vars)
		flen = len(p.funcs)
		dlen = len(defers)
		kws  = fract.None
	)
	(&obj.Block{
		Try: func() {
			for {
				p.i++
				tks := p.Tks[p.i]
				if tks[0].T == fract.End { // Block is ended.
					break
				} else if tks[0].T == fract.Catch { // Catch.
					p.skipBlock(false)
					break
				}
				if kws = p.process(tks); kws != fract.None {
					p.skipBlock(false)
				}
			}
			fract.TryCount--
			p.vars = p.vars[:vlen]
			p.funcs = p.funcs[:flen]
			for index := len(defers) - 1; index >= dlen; index-- {
				defers[index].call()
			}
			defers = defers[:dlen]
		},
		Catch: func(e obj.Exception) {
			p.loopCount = 0
			fract.TryCount--
			p.vars = p.vars[:vlen]
			p.funcs = p.funcs[:flen]
			defers = defers[:dlen]
			c := 0
			for {
				p.i++
				tks = p.Tks[p.i]
				if tks[0].T == fract.End {
					c--
					if c < 0 {
						break
					}
				} else if IsBlock(tks) {
					c++
				}
				if c > 0 {
					continue
				}
				if tks[0].T == fract.Catch {
					break
				}
			}
			// Ended block.
			if c < 0 {
				return
			}
			// Catch block.
			if len(tks) > 1 {
				fract.Error(tks[1], "Invalid syntax!")
			}
			for {
				p.i++
				tks := p.Tks[p.i]
				if tks[0].T == fract.End { // Block is ended.
					break
				}
				if kws = p.process(tks); kws != fract.None {
					p.skipBlock(false)
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
