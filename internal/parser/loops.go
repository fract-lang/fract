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
	bi := findBlock(tks)
	btks, tks := p.getBlock(tks[bi:]), tks[1:bi]
	flen := len(p.funcs)
	brk := false
	kws := fract.None
	//*************
	//    WHILE
	//*************
	if len(tks) == 0 || len(tks) >= 1 {
		if len(tks) == 0 || len(tks) == 1 || len(tks) >= 1 && tks[1].T != fract.In && tks[1].T != fract.Comma {
			vlen := len(p.vars)
			// Infinity loop.
			if len(tks) == 0 {
			infinity:
				for _, tks := range btks {
					kws = p.process(tks)
					if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return.
						return prockws(kws)
					} else if kws == fract.LOOPContinue { // Continue loop.
						break
					}
				}
				// Remove temporary variables.
				p.vars = p.vars[:vlen]
				// Remove temporary functions.
				p.funcs = p.funcs[:flen]
				goto infinity
			}
		while:
			// Interpret/skip block.
			c := p.procCondition(tks)
			for _, tks := range btks {
				// Condition is true?
				if c == "true" {
					kws = p.process(tks)
					if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return.
						brk = true
						break
					} else if kws == fract.LOOPContinue { // Continue loop.
						break
					}
				} else {
					brk = true
					break
				}
			}
			// Remove temporary variables.
			p.vars = p.vars[:vlen]
			// Remove temporary functions.
			p.funcs = p.funcs[:flen]
			c = p.procCondition(tks)
			if brk || c != "true" {
				return prockws(kws)
			}
			goto while
		}
	}

	//*************
	//   FOREACH
	//*************
	nametk := tks[0]
	// Name is not name?
	if nametk.T != fract.Name {
		fract.IPanic(nametk, obj.SyntaxPanic, "This is not a valid name!")
	}
	if ln := p.definedName(nametk); ln != -1 {
		fract.IPanic(nametk, obj.NamePanic, "\""+nametk.V+"\" is already defined at line: "+fmt.Sprint(ln))
	}
	// Element name?
	ename := ""
	if tks[1].T == fract.Comma {
		if len(tks) < 3 || tks[2].T != fract.Name {
			fract.IPanic(tks[1], obj.SyntaxPanic, "Element name is not defined!")
		}
		if tks[2].V != "_" {
			ename = tks[2].V
			if ln := p.definedName(tks[2]); ln != -1 {
				fract.IPanic(tks[2], obj.NamePanic, "\""+ename+"\" is already defined at line: "+fmt.Sprint(ln))
			}
		}
		if len(tks)-3 == 0 {
			tks[2].Col += len(tks[2].V)
			fract.IPanic(tks[2], obj.SyntaxPanic, "Value is not given!")
		}
		tks = tks[2:]
	}
	if vtks, inTk := tks.Sub(2, len(tks)-2), tks[1]; vtks != nil {
		tks = *vtks
	} else {
		fract.IPanic(inTk, obj.SyntaxPanic, "Value is not given!")
	}
	v := p.procVal(tks)
	// Type is not array?
	if !v.Arr && v.D[0].T != obj.VStr {
		fract.IPanic(tks[0], obj.ValuePanic, "Foreach loop must defined array value!")
	}
	p.vars = append(p.vars,
		obj.Var{Name: nametk.V, V: obj.Value{D: []obj.Data{{D: "0", T: obj.VInt}}}},
		obj.Var{Name: ename, V: obj.Value{}},
	)
	vlen := len(p.vars)
	index := &p.vars[vlen-2]
	element := &p.vars[vlen-1]
	vars := p.vars
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
			element.V.D = []obj.Data{v.D[0]}
		} else {
			element.V.D = []obj.Data{{D: string(v.D[0].String()[0]), T: obj.VStr}}
		}
	}
	// Interpret block.
	for j := 0; j < length; {
		for _, tks := range btks {
			kws = p.process(tks)
			if kws == fract.LOOPBreak || kws == fract.FUNCReturn { // Break loop or return.
				brk = true
				break
			} else if kws == fract.LOOPContinue { // Continue loop.
				break
			}
		}
		// Remove temporary variables.
		p.vars = vars
		// Remove temporary functions.
		p.funcs = p.funcs[:flen]
		j++
		if brk || (v.Arr && j == len(v.D) || !v.Arr && j == len(v.D[0].String())) {
			break
		}
		if index.Name != "" {
			index.V.D = []obj.Data{{D: fmt.Sprint(j), T: obj.VInt}}
		}
		if element.Name != "" {
			if v.Arr {
				element.V.D = []obj.Data{v.D[j]}
			} else {
				element.V.D = []obj.Data{{D: string(v.D[0].String()[j]), T: obj.VStr}}
			}
		}
	}
	// Remove loop variables.
	p.vars = vars[:len(vars)-2]
	return prockws(kws)
}
