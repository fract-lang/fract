package parser

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// prockws returns return value of kwstate.
func prockws(kws uint8) uint8 {
	if kws != fract.FUNCReturn {
		return fract.None
	}
	return kws
}

// procLoop process loops and returns keyword state.
func (p *Parser) procLoop(tks obj.Tokens) uint8 {
	// Content is empty?
	if vtokens := tks.Sub(1, len(tks)-1); vtokens == nil {
		tks = nil
	} else {
		tks = *vtokens
	}
	flen := len(p.funcs)
	brk := false
	kws := fract.None
	iindex := p.i

	//*************
	//    WHILE
	//*************
	if tks == nil || len(tks) >= 1 {
		if tks == nil || len(tks) == 1 || len(tks) >= 1 && tks[1].T != fract.In && tks[1].T != fract.Comma {
			vlen := len(p.vars)
			/* Infinity loop. */
			if tks == nil {
				for {
					p.i++
					tks := p.Tks[p.i]
					if tks[0].T == fract.End { // Block is ended.
						// Remove temporary variables.
						p.vars = p.vars[:vlen]
						// Remove temporary functions.
						p.funcs = p.funcs[:flen]
						if brk {
							return prockws(kws)
						}
						p.i = iindex
						continue
					} else if tks[0].T == fract.Else { // Else block.
						if len(tks) > 1 {
							fract.Error(tks[0], "Else block is not take any arguments!")
						}
						p.skipBlock(false)
						p.i--
						continue
					}
					kws = p.process(tks)
					if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return?
						brk = true
						p.skipBlock(false)
						p.i--
					} else if kws == fract.LOOPContinue { // Continue loop?
						p.skipBlock(false)
						p.i--
					}
				}
			}

			/* Interpret/skip block. */
			ctks := tks
			c := p.procCondition(ctks)
			_else := c == "false"
			for {
				p.i++
				tks := p.Tks[p.i]

				if tks[0].T == fract.End { // Block is ended.
					// Remove temporary variables.
					p.vars = p.vars[:vlen]
					// Remove temporary functions.
					p.funcs = p.funcs[:flen]
					c = p.procCondition(ctks)
					if brk || c != "true" {
						return prockws(kws)
					}
					p.i = iindex
					continue
				} else if tks[0].T == fract.Else { // Else block.
					if len(tks) > 1 {
						fract.Error(tks[0], "Else block is not take any arguments!")
					}
					if c == "true" {
						p.skipBlock(false)
						p.i--
						continue
					}
					// Remove temporary variables.
					p.vars = p.vars[:vlen]
					// Remove temporary functions.
					p.funcs = p.funcs[:flen]
					if !_else {
						p.skipBlock(false)
						return prockws(kws)
					}
					for {
						p.i++
						tks = p.Tks[p.i]
						if tks[0].T == fract.End { // Block is ended.
							// Remove temporary variables.
							p.vars = p.vars[:vlen]
							// Remove temporary functions.
							p.funcs = p.funcs[:flen]
							return prockws(kws)
						}
						kws = p.process(tks)
						if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return?
							brk = true
							p.skipBlock(false)
							p.i--
						} else if kws == fract.LOOPContinue { // Continue loop?
							p.skipBlock(false)
							p.i--
						}
					}
				}
				// Condition is true?
				if c == "true" {
					kws = p.process(tks)
					if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return?
						brk = true
						p.skipBlock(false)
						p.i--
					} else if kws == fract.LOOPContinue { // Continue loop?
						p.skipBlock(false)
						p.i--
					}
				} else {
					if _else {
						p.skipBlock(true)
						continue
					}
					brk = true
					p.skipBlock(false)
					p.i--
				}
			}
		}
	}

	//*************
	//   FOREACH
	//*************
	nametk := tks[0]
	// Name is not name?
	if nametk.T != fract.Name {
		fract.Error(nametk, "This is not a valid name!")
	}
	if ln := p.definedName(nametk); ln != -1 {
		fract.Error(nametk, "\""+nametk.Val+"\" is already defined at line: "+fmt.Sprint(ln))
	}
	// Element name?
	ename := ""
	if tks[1].T == fract.Comma {
		if len(tks) < 3 || tks[2].T != fract.Name {
			fract.Error(tks[1], "Element name is not defined!")
		}
		if tks[2].Val != "_" {
			ename = tks[2].Val
			if ln := p.definedName(tks[2]); ln != -1 {
				fract.Error(tks[2], "\""+ename+"\" is already defined at line: "+fmt.Sprint(ln))
			}
		}
		if len(tks)-3 == 0 {
			tks[2].Col += len(tks[2].Val)
			fract.Error(tks[2], "Value is not defined!")
		}
		tks = tks[2:]
	}
	if vtks, inTk := tks.Sub(2, len(tks)-2), tks[1]; vtks != nil {
		tks = *vtks
	} else {
		fract.Error(inTk, "Value is not defined!")
	}
	v := p.procVal(tks)
	// Type is not array?
	if !v.Arr && v.D[0].T != obj.VStr {
		fract.Error(tks[0], "Foreach loop must defined array value!")
	}
	// Empty array?
	if v.Arr && len(v.D) == 0 || v.D[0].T == obj.VStr && v.D[0].D == "" {
		vlen := len(p.vars)
		for {
			p.i++
			tks := p.Tks[p.i]
			if tks[0].T == fract.End { // Block is ended.
				return kws
			} else if tks[0].T == fract.Else { // Else block.
				if len(tks) > 1 {
					fract.Error(tks[0], "Else block is not take any arguments!")
				}
				for {
					p.i++
					tks = p.Tks[p.i]
					if tks[0].T == fract.End { // Block is ended.
						// Remove temporary variables.
						p.vars = p.vars[:vlen]
						// Remove temporary functions.
						p.funcs = p.funcs[:flen]
						return prockws(kws)
					}
					kws = p.process(tks)
					if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return?
						brk = true
						p.skipBlock(false)
						p.i--
					} else if kws == fract.LOOPContinue { // Continue loop?
						p.skipBlock(false)
						p.i--
					}
				}
			}
			p.skipBlock(true)
		}
	}
	p.vars = append(p.vars,
		obj.Var{Name: nametk.Val, Val: obj.Value{D: []obj.Data{{D: "0", T: obj.VInt}}}},
		obj.Var{Name: ename, Val: obj.Value{}},
	)
	vlen := len(p.vars)
	index := &p.vars[vlen-2]
	element := &p.vars[vlen-1]
	if index.Name == "_" {
		index.Name = ""
	}
	var length int
	if v.Arr {
		length = len(v.D)
	} else {
		length = len(v.D[0].String())
	}
	if element.Name != "" {
		if v.Arr {
			element.Val.D = []obj.Data{v.D[0]}
		} else {
			element.Val.D = []obj.Data{{D: string(v.D[0].String()[0]), T: obj.VStr}}
		}
	}
	// Interpret block.
	for j := 0; j < length; {
		p.i++
		tks := p.Tks[p.i]
		if tks[0].T == fract.End { // Block is ended.
			// Remove temporary variables.
			p.vars = p.vars[:vlen]
			// Remove temporary functions.
			p.funcs = p.funcs[:flen]
			j++
			if brk || (v.Arr && j == len(v.D) || !v.Arr && j == len(v.D[0].String())) {
				break
			}
			p.i = iindex
			if index.Name != "" {
				index.Val.D = []obj.Data{{D: fmt.Sprint(j), T: obj.VInt}}
			}
			if element.Name != "" {
				if v.Arr {
					element.Val.D = []obj.Data{v.D[j]}
				} else {
					element.Val.D = []obj.Data{{D: string(v.D[0].String()[j]), T: obj.VStr}}
				}
			}
			continue
		} else if tks[0].T == fract.Else { // Else block.
			if len(tks) > 1 {
				fract.Error(tks[0], "Else block is not take any arguments!")
			}
			p.skipBlock(false)
			p.i--
			continue
		}
		kws = p.process(tks)
		if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return?
			brk = true
			p.skipBlock(false)
			p.i--
		} else if kws == fract.LOOPContinue { // Continue next?
			p.skipBlock(false)
			p.i--
		}
	}
	// Remove loop variables.
	p.vars = p.vars[:vlen-2]
	return prockws(kws)
}
