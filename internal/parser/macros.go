package parser

import (
	"fmt"
	"runtime"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

func (p *Parser) processMacroIf(tks obj.Tokens) uint8 {
	tlen := len(tks)
	ctks := tks.Sub(1, tlen-1)
	// Condition is empty?
	if ctks == nil {
		first := tks[0]
		fract.Errorc(first.F, first.Ln, first.Col+len(first.Val), "Condition is empty!")
	}
	vars := p.vars
	funcs := p.funcs
	p.vars = append([]obj.Var{
		{
			Name: "OS",
			Val:  obj.Value{D: []obj.Data{{D: runtime.GOOS, T: obj.VString}}},
		},
		{
			Name: "ARCH",
			Val: obj.Value{
				D: []obj.Data{{D: runtime.GOARCH, T: obj.VString}},
			},
		},
	}, p.macroDefs...)
	state := p.processCondition(*ctks)
	kws := fract.None
	/* Interpret/skip block. */
	for {
		p.i++
		tks := p.Tks[p.i]
		first := tks[0]
		if first.T == fract.Macro {
			tks := tks[1:]
			first = tks[0]
			if first.T == fract.End { // Block is ended.
				goto ret
			} else if first.T == fract.ElseIf { // Else if block.
				tlen = len(tks)
				ctks := tks.Sub(1, tlen-1)
				// Condition is empty?
				if ctks == nil {
					first := tks[0]
					fract.Errorc(first.F, first.Ln, first.Col+len(first.Val), "Condition is empty!")
				}
				if state == "true" {
					p.skipBlock(false)
					goto ret
				}
				state = p.processCondition(*ctks)
				// Interpret/skip block.
				for {
					p.i++
					tks := p.Tks[p.i]
					first := tks[0]
					if first.T == fract.Macro {
						tks := tks[1:]
						first = tks[0]
						if first.T == fract.End { // Block is ended.
							goto ret
						} else if first.T == fract.If { // If block.
							if state == "true" && kws == fract.None {
								p.processMacroIf(tks)
							} else {
								p.skipBlock(true)
							}
							continue
						} else if first.T == fract.ElseIf || first.T == fract.Else { // Else if or else block.
							p.i--
							break
						}
					}
					// Condition is true?
					if state == "true" && kws == fract.None {
						p.vars, vars = vars, p.vars
						kws = p.processTokens(tks)
						p.vars, vars = vars, p.vars
						if kws != fract.None {
							p.skipBlock(false)
						}
					} else {
						p.skipBlock(true)
					}
				}
				if state == "true" {
					p.skipBlock(false)
					goto ret
				}
				continue
			} else if first.T == fract.Else { // Else block.
				if len(tks) > 1 {
					fract.Error(first, "Else block is not take any arguments!")
				}
				if state == "true" {
					p.skipBlock(false)
					goto ret
				}
				/* Interpret/skip block. */
				for {
					p.i++
					tks := p.Tks[p.i]
					first := tks[0]
					if first.T == fract.Macro {
						tks = tks[1:]
						first = tks[0]
						if first.T == fract.End { // Block is ended.
							goto ret
						} else if first.T == fract.If { // If block.
							if kws == fract.None {
								p.processMacroIf(tks)
							} else {
								p.skipBlock(true)
							}
							continue
						}
					}
					// Condition is true?
					if kws == fract.None {
						p.vars, vars = vars, p.vars
						kws = p.processTokens(tks)
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
			kws = p.processTokens(tks)
			p.vars, vars = vars, p.vars
			if kws != fract.None {
				p.skipBlock(false)
			}
		} else {
			p.skipBlock(true)
		}
	}
ret:
	p.vars = vars
	p.funcs = funcs
	return kws
}

func (p *Parser) processMacroDefine(tks []obj.Token) obj.Var {
	if len(tks) < 2 {
		fract.Error(tks[0], "Define name is not defined!")
	}
	name := tks[1]
	if name.T != fract.Name {
		fract.Error(name, "Invalid name!")
	}
	// Exists name.
	for _, macro := range p.macroDefs {
		if macro.Name == name.Val {
			fract.Error(name, "This macro is already defined in this name at line: "+fmt.Sprint(macro.Ln))
		}
	}
	macro := obj.Var{
		Name: name.Val,
		Ln:   name.Ln,
	}
	if len(tks) > 2 {
		vars := p.vars
		macro.Val = p.processValue(tks[2:])
		p.vars = vars
	} else {
		macro.Val.D = []obj.Data{{
			D: "false",
			T: obj.VBoolean,
		}}
	}
	return macro
}

// processMacro process macros and returns keyword state.
func (p *Parser) processMacro(tks []obj.Token) uint8 {
	tks = tks[1:]
	switch tks[0].T {
	case fract.If:
		return p.processMacroIf(tks)
	case fract.Name:
		switch tks[0].Val {
		case "define": // Macro variable.
			p.macroDefs = append(p.macroDefs, p.processMacroDefine(tks))
		case "pragma":
			if len(tks) != 2 || tks[1].T != fract.Name {
				fract.Error(tks[0], "Invalid pragma syntax!")
			}
			switch tks[1].Val {
			case "enofi":
				if p.loopCount == -1 {
					p.loopCount = 0
				}
			default:
				fract.Error(tks[1], "Invalid pragma!")
			}
		default:
			fract.Error(tks[0], "Invalid macro!")
		}
	default:
		fract.Error(tks[0], "Invalid macro!")
	}
	return fract.None
}
