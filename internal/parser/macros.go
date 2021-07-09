package parser

import (
	"runtime"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

func (p *Parser) procMacroIf(tks obj.Tokens) uint8 {
	tlen := len(tks)
	ctks := tks.Sub(1, tlen-1)
	// Condition is empty?
	if ctks == nil {
		first := tks[0]
		fract.IPanicC(first.F, first.Ln, first.Col+len(first.V), obj.SyntaxPanic, "Condition is not given!")
	}
	vars := p.vars
	funcs := p.funcs
	p.vars = []obj.Var{{
		Name: "OS",
		V:    obj.Value{D: []obj.Data{{D: runtime.GOOS, T: obj.VStr}}},
	}, {
		Name: "ARCH",
		V:    obj.Value{D: []obj.Data{{D: runtime.GOARCH, T: obj.VStr}}},
	}}
	state := p.procCondition(*ctks)
	kws := fract.None
	// Interpret/skip block.
	for {
		p.i++
		tks := p.Tks[p.i]
		if tks[0].T == fract.Macro {
			tks := tks[1:]
			switch tks[0].T {
			case fract.End: // Block is ended.
				goto end
			case fract.ElseIf: // Else if block.
				tlen = len(tks)
				ctks := tks.Sub(1, tlen-1)
				// Condition is empty?
				if ctks == nil {
					first := tks[0]
					fract.IPanicC(first.F, first.Ln, first.Col+len(first.V), obj.ValuePanic, "Condition is empty!")
				}
				if state == "true" {
					p.skipBlock(false)
					goto end
				}
				state = p.procCondition(*ctks)
				// Interpret/skip block.
				for {
					p.i++
					tks := p.Tks[p.i]
					if tks[0].T == fract.Macro {
						tks := tks[1:]
						switch tks[0].T {
						case fract.End: // Block is ended.
							goto end
						case fract.If: // If block.
							if state == "true" && kws == fract.None {
								p.procMacroIf(tks)
							} else {
								p.skipBlock(true)
							}
							continue
						case fract.ElseIf, fract.Else: // Else if or else block.
							p.i--
							goto elifend
						}
					}
					// Condition is true?
					if state == "true" && kws == fract.None {
						p.vars, vars = vars, p.vars
						kws = p.process(tks)
						p.vars, vars = vars, p.vars
						if kws != fract.None {
							p.skipBlock(false)
						}
					} else {
						p.skipBlock(true)
					}
				}
			elifend:
				if state == "true" {
					p.skipBlock(false)
					goto end
				}
				continue
			case fract.Else: // Else block.
				if len(tks) > 1 {
					fract.IPanic(tks[0], obj.SyntaxPanic, "Else block is not take any arguments!")
				}
				if state == "true" {
					p.skipBlock(false)
					goto end
				}
				/* Interpret/skip block. */
				for {
					p.i++
					tks := p.Tks[p.i]
					first := tks[0]
					if first.T == fract.Macro {
						tks = tks[1:]
						switch tks[0].T {
						case fract.End: // Block is ended.
							goto end
						case fract.If: // If block.
							if kws == fract.None {
								p.procMacroIf(tks)
							} else {
								p.skipBlock(true)
							}
							continue
						}
					}
					// Condition is true?
					if kws == fract.None {
						p.vars, vars = vars, p.vars
						kws = p.process(tks)
						p.vars, vars = vars, p.vars
						if kws != fract.None {
							p.skipBlock(false)
						}
					}
				}
			}
		}
		// Condition is true?
		if state == "true" && kws == fract.None {
			p.vars, vars = vars, p.vars
			kws = p.process(tks)
			p.vars, vars = vars, p.vars
			if kws != fract.None {
				p.skipBlock(false)
			}
		} else {
			p.skipBlock(true)
		}
	}
end:
	p.vars = vars
	p.funcs = funcs
	return kws
}

// procMacro process macros and returns keyword state.
func (p *Parser) procMacro(tks []obj.Token) uint8 {
	tks = tks[1:]
	switch tks[0].T {
	case fract.If:
		return p.procMacroIf(tks)
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
