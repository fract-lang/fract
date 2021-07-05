package parser

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// Compare arithmetic values.
func compVals(opr string, d0, d1 obj.Data) bool {
	if d0.T != d1.T && (d0.T == obj.VStr || d1.T == obj.VStr) {
		return false
	}
	switch opr {
	case "==": // Equals.
		if (d0.T == obj.VStr && d0.D != d1.D) ||
			(d0.T != obj.VStr && arithmetic.Value(d0.String()) != arithmetic.Value(d1.String())) {
			return false
		}
	case "<>": // Not equals.
		if (d0.T == obj.VStr && d0.D == d1.D) ||
			(d0.T != obj.VStr && arithmetic.Value(d0.String()) == arithmetic.Value(d1.String())) {
			return false
		}
	case ">": // Greater.
		if (d0.T == obj.VStr && d0.String() <= d1.String()) ||
			(d0.T != obj.VStr && arithmetic.Value(d0.String()) <= arithmetic.Value(d1.String())) {
			return false
		}
	case "<": // Less.
		if (d0.T == obj.VStr && d0.String() >= d1.String()) ||
			(d0.T != obj.VStr && arithmetic.Value(d0.String()) >= arithmetic.Value(d1.String())) {
			return false
		}
	case ">=": // Greater or equals.
		if (d0.T == obj.VStr && d0.String() < d1.String()) ||
			(d0.T != obj.VStr && arithmetic.Value(d0.String()) < arithmetic.Value(d1.String())) {
			return false
		}
	case "<=": // Less or equals.
		if (d0.T == obj.VStr && d0.String() > d1.String()) ||
			(d0.T != obj.VStr && arithmetic.Value(d0.String()) > arithmetic.Value(d1.String())) {
			return false
		}
	}
	return true
}

// Compare values.
func comp(v0, v1 obj.Value, opr obj.Token) bool {
	// In.
	if opr.Val == "in" {
		if !v1.Arr && v1.D[0].T != obj.VStr {
			fract.Error(opr, "Value is not enumerable!")
		}
		if v1.Arr {
			dt := v0.String()
			for _, d := range v1.D {
				if strings.Contains(d.String(), dt) {
					return true
				}
			}
		} else { // String.
			if v0.Arr {
				dt := v1.D[0].String()
				for _, d := range v0.D {
					if d.T != obj.VStr {
						fract.Error(opr, "All datas is not string!")
					}
					if strings.Contains(dt, d.String()) {
						return true
					}
				}
			} else {
				if v1.D[0].T != obj.VStr {
					fract.Error(opr, "All datas is not string!")
				}
				if strings.Contains(v1.D[0].String(), v1.D[0].String()) {
					return true
				}
			}
		}
		return false
	}
	// String comparison.
	if !v0.Arr || !v1.Arr {
		d0, d1 := v0.D[0], v1.D[0]
		if (d0.T == obj.VStr && d1.T != obj.VStr) || (d0.T != obj.VStr && d1.T == obj.VStr) {
			fract.Error(opr, "The in keyword should use with string or enumerable data types!")
		}
		return compVals(opr.Val, d0, d1)
	}
	// Array comparison.
	if v0.Arr || v1.Arr {
		if (v0.Arr && !v1.Arr) || (!v0.Arr && v1.Arr) {
			return false
		}
		if len(v0.D) != len(v1.D) {
			return opr.Val == "<>"
		}
		for i, d := range v0.D {
			if !compVals(opr.Val, d, v1.D[i]) {
				return false
			}
		}
		return true
	}
	// Single value comparison.
	return compVals(opr.Val, v0.D[0], v1.D[0])
}

// procCondition returns condition result.
func (p *Parser) procCondition(tks obj.Tokens) string {
	p.procRange(&tks)
	T := obj.Value{D: []obj.Data{{D: "true"}}}
	// Process condition.
	ors := conditionalProcesses(tks, "||")
	for _, or := range *ors {
		// Decompose and conditions.
		ands := conditionalProcesses(or, "&&")
		// Is and long statement?
		if len(*ands) > 1 {
			for _, and := range *ands {
				i, opr := findConditionOpr(and)
				// Operator is not found?
				if i == -1 {
					opr.Val = "=="
					if comp(p.procVal(and), T, opr) {
						return "true"
					}
					continue
				}
				// Operator is first or last?
				if i == 0 {
					fract.Error(and[0], "Comparison values are missing!")
				} else if i == len(and)-1 {
					fract.Error(and[len(and)-1], "Comparison values are missing!")
				}
				if !comp(
					p.procVal(*and.Sub(0, i)), p.procVal(*and.Sub(i+1, len(and)-i-1)), opr) {
					return "false"
				}
			}
			return "true"
		}
		i, opr := findConditionOpr(or)
		// Operator is not found?
		if i == -1 {
			opr.Val = "=="
			if comp(p.procVal(or), T, opr) {
				return "true"
			}
			continue
		}
		// Operator is first or last?
		if i == 0 {
			fract.Error(or[0], "Comparison values are missing!")
		} else if i == len(or)-1 {
			fract.Error(or[len(or)-1], "Comparison values are missing!")
		}
		if comp(p.procVal(*or.Sub(0, i)), p.procVal(*or.Sub(i+1, len(or)-i-1)), opr) {
			return "true"
		}
	}
	return "false"
}

// Get string arithmetic compatible data.
func arith(tks obj.Token, d obj.Data) string {
	ret := d.String()
	switch d.T {
	case obj.VFunc:
		fract.Error(tks, "\""+ret+"\" is not compatible with arithmetic processes!")
	}
	return ret
}

// process instance for solver.
type process struct {
	f   obj.Token // First value of process.
	fv  obj.Value // Value instance of first value.
	s   obj.Token // Second value of process.
	sv  obj.Value // Value instance of second value.
	opr obj.Token // Operator of process.
}

// Tokenize array.
func tokenizeArray(dts []obj.Data) obj.Tokens {
	tks := obj.Tokens{obj.Token{Val: "[", T: fract.Brace}}
	for _, d := range dts {
		switch d.T {
		case obj.VArray:
			tks = append(tks, tokenizeArray(d.D.([]obj.Data))...)
		default:
			tks = append(tks, obj.Token{Val: d.String(), T: fract.Value})
		}
		tks = append(tks, obj.Token{Val: ",", T: fract.Comma})
	}
	tks[len(tks)-1] = obj.Token{Val: "]", T: fract.Brace}
	return tks
}

// procRange by value processor principles.
func (p *Parser) procRange(tks *obj.Tokens) {
	for {
		rg, i := decomposeBrace(tks, "(", ")", true)
		/* Parentheses are not found! */
		if i == -1 {
			return
		}
		val := p.procVal(rg)
		if val.Arr {
			tks.Ins(i, tokenizeArray(val.D)...)
		} else {
			if val.D[0].T == obj.VStr {
				tks.Ins(i, obj.Token{Val: "'" + val.D[0].String() + "'", T: fract.Value})
			} else {
				tks.Ins(i, obj.Token{
					Val: val.D[0].String(),
					T:   fract.Value,
					//! Add another fields for panic.
					Ln:  rg[0].Ln,
					Col: rg[0].Col,
					F:   rg[0].F,
				})
			}
		}
	}
}

// solve process.
func solve(opr obj.Token, a, b float64) float64 {
	var r float64
	if opr.Val == "\\" || opr.Val == "\\\\" { // Divide with bigger.
		if opr.Val == "\\" {
			opr.Val = "/"
		} else {
			opr.Val = "//"
		}
		if a < b {
			cache := a
			a = b
			b = cache
		}
	}
	switch opr.Val {
	case "+": // Addition.
		r = a + b
	case "-": // Subtraction.
		r = a - b
	case "*": // Multiply.
		r = a * b
	case "/", "//": // Division.
		if a == 0 || b == 0 {
			fract.Error(opr, "Divide by zero!")
		}
		r = a / b
		if opr.Val == "//" {
			r = math.RoundToEven(r)
		}
	case "|": // Binary or.
		r = float64(int(a) | int(b))
	case "&": // Binary and.
		r = float64(int(a) & int(b))
	case "^": // Bitwise exclusive or.
		r = float64(int(a) ^ int(b))
	case "**": // Exponentiation.
		r = math.Pow(a, b)
	case "%": // Mod.
		r = math.Mod(a, b)
	case "<<": // Left shift.
		if b < 0 {
			fract.Error(opr, "Shifter is cannot should be negative!")
		}
		r = float64(int(a) << int(b))
	case ">>": // Right shift.
		if b < 0 {
			fract.Error(opr, "Shifter is cannot should be negative!")
		}
		r = float64(int(a) >> int(b))
	default:
		fract.Error(opr, "Operator is invalid!")
	}
	return r
}

// Check data and set ready.
func readyData(p process, d obj.Data) obj.Data {
	if p.fv.D[0].T == obj.VStr || p.sv.D[0].T == obj.VStr {
		d.T = obj.VStr
	} else if p.opr.Val == "/" || p.opr.Val == "\\" ||
		p.fv.D[0].T == obj.VFloat || p.sv.D[0].T == obj.VFloat {
		d.T = obj.VFloat
		return d
	}
	return d
}

// solveProc solve arithmetic process.
func solveProc(p process) obj.Value {
	v := obj.Value{D: []obj.Data{{}}}
	// String?
	if (len(p.fv.D) != 0 && p.fv.D[0].T == obj.VStr) || (len(p.sv.D) != 0 && p.sv.D[0].T == obj.VStr) {
		if p.fv.D[0].T == p.sv.D[0].T { // Both string?
			v.D[0].T = obj.VStr
			switch p.opr.Val {
			case "+":
				v.D[0].D = p.fv.D[0].String() + p.sv.D[0].String()
			case "-":
				flen := len(p.fv.D[0].String())
				slen := len(p.sv.D[0].String())
				if flen == 0 || slen == 0 {
					v.D[0].D = ""
					break
				}
				if flen == 1 && slen > 1 {
					r, _ := strconv.ParseInt(p.fv.D[0].String(), 10, 32)
					fr := rune(r)
					for _, r := range p.sv.D[0].String() {
						v.D[0].D = v.D[0].String() + string(fr-r)
					}
				} else if slen == 1 && flen > 1 {
					r, _ := strconv.ParseInt(p.sv.D[0].String(), 10, 32)
					fr := rune(r)
					for _, r := range p.fv.D[0].String() {
						v.D[0].D = v.D[0].String() + string(fr-r)
					}
				} else {
					for i, r := range p.fv.D[0].String() {
						v.D[0].D = v.D[0].String() + string(r-rune(p.sv.D[0].String()[i]))
					}
				}
			default:
				fract.Error(p.opr, "This operator is not defined for string types!")
			}
			return v
		}

		v.D[0].T = obj.VStr
		if p.fv.D[0].T == obj.VStr {
			if p.sv.Arr {
				if len(p.sv.D) == 0 {
					v.D = p.fv.D
					return v
				}
				if len(p.fv.D[0].String()) != len(p.sv.D) && (len(p.fv.D[0].String()) != 1 && len(p.sv.D) != 1) {
					fract.Error(p.s, "Array element count is not one or equals to first array!")
				}
				if strings.Contains(p.sv.D[0].String(), ".") {
					fract.Error(p.s, "Only string and integer values cannot concatenate string values!")
				}
				r, _ := strconv.ParseInt(p.sv.D[0].String(), 10, 64)
				rn := rune(r)
				var sb strings.Builder
				for _, r := range p.fv.D[0].String() {
					switch p.opr.Val {
					case "+":
						sb.WriteByte(byte(r + rn))
					case "-":
						sb.WriteByte(byte(r - rn))
					default:
						fract.Error(p.opr, "This operator is not defined for string types!")
					}
				}
				v.D[0].D = sb.String()
			} else {
				if p.sv.D[0].T != obj.VInt {
					fract.Error(p.s, "Only string and integer values cannot concatenate string values!")
				}
				var sb strings.Builder
				rs, _ := strconv.ParseInt(p.sv.D[0].String(), 10, 64)
				rn := rune(rs)
				for _, r := range p.fv.D[0].String() {
					switch p.opr.Val {
					case "+":
						sb.WriteByte(byte(r + rn))
					case "-":
						sb.WriteByte(byte(r - rn))
					default:
						fract.Error(p.opr, "This operator is not defined for string types!")
					}
				}
				v.D[0].D = sb.String()
			}
		} else {
			if p.fv.Arr {
				if len(p.fv.D) == 0 {
					v.D = p.sv.D
					return v
				}
				if len(p.fv.D[0].String()) != len(p.sv.D) && (len(p.fv.D[0].String()) != 1 && len(p.sv.D) != 1) {
					fract.Error(p.s, "Array element count is not one or equals to first array!")
				}
				if strings.Contains(p.fv.D[0].String(), ".") {
					fract.Error(p.s, "Only string and integer values cannot concatenate string values!")
				}
				rs, _ := strconv.ParseInt(p.fv.D[0].String(), 10, 64)
				rn := rune(rs)
				var sb strings.Builder
				for _, r := range p.sv.D[0].String() {
					switch p.opr.Val {
					case "+":
						sb.WriteByte(byte(r + rn))
					case "-":
						sb.WriteByte(byte(r - rn))
					default:
						fract.Error(p.opr, "This operator is not defined for string types!")
					}
				}
				v.D[0].D = sb.String()
			} else {
				if p.fv.D[0].T != obj.VInt {
					fract.Error(p.f, "Only string and integer values cannot concatenate string values!")
				}
				var sb strings.Builder
				rs, _ := strconv.ParseInt(p.fv.D[0].String(), 10, 64)
				rn := rune(rs)
				for _, r := range p.sv.D[0].String() {
					switch p.opr.Val {
					case "+":
						sb.WriteByte(byte(r + rn))
					case "-":
						sb.WriteByte(byte(r - rn))
					default:
						fract.Error(p.opr, "This operator is not defined for string types!")
					}
				}
				v.D[0].D = sb.String()
			}
		}
		return v
	}

	if p.fv.Arr && p.sv.Arr {
		v.Arr = true
		if len(p.fv.D) == 0 {
			v.D = p.sv.D
			return v
		} else if len(p.sv.D) == 0 {
			v.D = p.fv.D
			return v
		}
		if len(p.fv.D) != len(p.sv.D) && (len(p.fv.D) != 1 && len(p.sv.D) != 1) {
			fract.Error(p.s, "Array element count is not one or equals to first array!")
		}
		if len(p.fv.D) == 1 || len(p.sv.D) == 1 {
			f, s := p.fv, p.sv
			if len(f.D) != 1 {
				f, s = s, f
			}
			ar := arithmetic.Value(arith(p.opr, f.D[0]))
			for i, d := range s.D {
				if d.T == obj.VArray {
					s.D[i] = readyData(p, obj.Data{
						D: solveProc(process{
							f:  p.f,
							fv: s,
							s:  p.s,
							sv: obj.Value{
								D:   d.D.([]obj.Data),
								Arr: true,
							},
							opr: p.opr,
						}).D,
						T: obj.VArray,
					})
				} else {
					s.D[i] = readyData(p, obj.Data{D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, ar, arithmetic.Value(arith(p.opr, d))))})
				}
			}
			v.D = s.D
		} else {
			for i, f := range p.fv.D {
				s := p.sv.D[i]
				if f.T == obj.VArray || s.T == obj.VArray {
					proc := process{f: p.f, s: p.s, opr: p.opr}
					if f.T == obj.VArray {
						proc.fv = obj.Value{D: f.D.([]obj.Data), Arr: true}
					} else {
						proc.fv = obj.Value{D: []obj.Data{f}}
					}
					if s.T == obj.VArray {
						proc.sv = obj.Value{D: s.D.([]obj.Data), Arr: true}
					} else {
						proc.sv = obj.Value{D: []obj.Data{s}}
					}
					p.fv.D[i] = readyData(p, obj.Data{D: solveProc(proc).D, T: obj.VArray})
				} else {
					p.fv.D[i] = readyData(p,
						obj.Data{
							D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, arithmetic.Value(arith(p.opr, f)), arithmetic.Value(s.String()))),
						})
				}
			}
			v.D = p.fv.D
		}
	} else if p.fv.Arr || p.sv.Arr {
		v.Arr = true
		if len(p.fv.D) == 0 {
			v.D = p.sv.D
			return v
		} else if len(p.sv.D) == 0 {
			v.D = p.fv.D
			return v
		}
		f, s := p.fv, p.sv
		if !f.Arr {
			f, s = s, f
		}
		ar := arithmetic.Value(arith(p.opr, s.D[0]))
		for i, d := range f.D {
			if d.T == obj.VArray {
				f.D[i] = readyData(p, obj.Data{
					D: solveProc(process{
						f:   p.f,
						fv:  s,
						s:   p.s,
						sv:  obj.Value{D: d.D.([]obj.Data), Arr: true},
						opr: p.opr,
					}).D,
					T: obj.VArray,
				})
			} else {
				f.D[i] = readyData(p,
					obj.Data{
						D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, arithmetic.Value(arith(p.opr, d)), ar)),
					})
			}
		}
		v.D = f.D
	} else {
		if len(p.fv.D) == 0 {
			p.fv.D = []obj.Data{{D: "0"}}
		}
		v.D[0] = readyData(p,
			obj.Data{
				D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, arithmetic.Value(arith(p.opr, p.fv.D[0])), arithmetic.Value(arith(p.opr, p.sv.D[0])))),
			})
	}
	return v
}

// applyMinus operator.
func applyMinus(minus bool, v obj.Value) obj.Value {
	if !minus {
		return v
	}
	val := obj.Value{
		Arr: v.Arr,
		D:   append([]obj.Data{}, v.D...),
	}
	if val.Arr {
		for i, d := range val.D {
			if d.T == obj.VBool || d.T == obj.VFloat || d.T == obj.VInt {
				val.D[i].D = fmt.Sprintf(fract.FloatFormat, -arithmetic.Value(d.String()))
			}
		}
		return val
	}
	if d := val.D[0]; d.T == obj.VBool || d.T == obj.VFloat || d.T == obj.VInt {
		val.D[0].D = fmt.Sprintf(fract.FloatFormat, -arithmetic.Value(d.String()))
	}
	return val
}

// Process value part.
func (p *Parser) procValPart(fst bool, opr *process, tks *obj.Tokens, pos int) int {
	var (
		tk = opr.f
		r  = &opr.fv
	)
	if !fst {
		tk = opr.s
		r = &opr.sv
	}
	minus := tk.T == fract.Name && tk.Val[0] == '-'
	if tk.T == fract.Name {
		if pos < len(*tks)-1 {
			next := (*tks)[pos+1]
			// Array?
			if next.T == fract.Brace {
				switch next.Val {
				case "[":
					vi, t, src := p.defByName(tk)
					if vi == -1 || t != 'v' {
						fract.Error(tk, "Variable is not defined in this name: "+tk.Val)
					}
					// Find close bracket.
					ci := pos
					bc := 0
					for ; ci < len(*tks); ci++ {
						tk := (*tks)[ci]
						if tk.T == fract.Brace {
							if tk.Val == "[" {
								bc++
							} else if tk.Val == "]" {
								bc--
								if bc == 0 {
									break
								}
							}
						}
					}
					vtks := tks.Sub(pos+2, ci-pos-2)
					// Index value is empty?
					if vtks == nil {
						fract.Error(tk, "Index is not defined!")
					}
					v := src.vars[vi]
					if !v.Val.Arr && v.Val.D[0].T != obj.VStr {
						fract.Error((*tks)[pos], "Index accessor is cannot used with non-array variables!")
					}
					val := p.procVal(*vtks)
					i := indexes(v.Val, val, tk)
					tks.Rem(pos+1, ci-pos)
					var d []obj.Data
					if !v.Val.Arr {
						d = append(d, obj.Data{D: "", T: obj.VStr})
					}
					for _, pos := range i {
						if v.Val.Arr {
							d = append(d, v.Val.D[pos])
						} else {
							if v.Val.D[0].T == obj.VStr {
								d[0].D = d[0].String() + string(v.Val.D[0].String()[pos])
							} else {
								d[0].D = d[0].String() + fmt.Sprint(v.Val.D[0].String()[pos])
							}
						}
					}
					r.Arr = len(i) > 1 && d[0].T != obj.VStr || d[0].T == obj.VArray
					r.D = d
					*r = applyMinus(minus, *r)
					return 0
				case "(":
					// Find close parentheses.
					ci := pos + 1
					bc := 0
					for ; ci < len(*tks); ci++ {
						tk := (*tks)[ci]
						if tk.T == fract.Brace {
							if tk.Val == "(" {
								bc++
							} else if tk.Val == ")" {
								bc--
								if bc == 0 {
									break
								}
							}
						}
					}
					ci++
					v := p.funcCall((*tks)[pos:ci])
					if !opr.fv.Arr && v.D == nil {
						fract.Error(tk, "Function is not return any value!")
					}
					tks.Rem(pos+1, ci-pos-1)
					*r = applyMinus(minus, v)
					return 0
				}
			}
		}
		vi, t, src := p.defByName(tk)
		if vi == -1 {
			fract.Error(tk, "Variable is not defined in this name: "+tk.Val)
		}
		switch t {
		case 'f':
			*r = obj.Value{
				D: []obj.Data{{D: src.funcs[vi], T: obj.VFunc}},
			}
		case 'v':
			v := src.vars[vi]
			val := v.Val
			if !v.Mut { //! Immutability.
				val.D = append(make([]obj.Data, 0), v.Val.D...)
			}
			*r = applyMinus(minus, val)
		}
		return 0
	} else if tk.T == fract.Brace {
		switch tk.Val {
		case "}":
			// Find open bracket.
			bc := 1
			oi := pos - 1
			for ; oi >= 0; oi-- {
				tk := (*tks)[oi]
				if tk.T == fract.Brace {
					if tk.Val == "}" {
						bc++
					} else if tk.Val == "{" {
						bc--
						if bc == 0 {
							break
						}
					}
				}
			}
			// Finished?
			if oi == 0 || (*tks)[oi-1].T != fract.Name {
				r.Arr = true
				r.D = p.procEnumerableVal(*tks.Sub(oi, pos-oi+1)).D
				*r = applyMinus(minus, *r)
				tks.Rem(oi, pos-oi)
				return pos - oi
			}
			endtk := (*tks)[oi-1]
			vi, t, src := p.defByName(tk)
			if vi == -1 || t != 'v' {
				fract.Error(endtk, "Variable is not defined in this name: "+endtk.Val)
			}
			vtks := tks.Sub(oi+1, pos-oi-1)
			// Index value is empty?
			if vtks == nil {
				fract.Error(endtk, "Index is not defined!")
			}
			v := src.vars[vi]
			if !v.Val.Arr && v.Val.D[0].T != obj.VStr {
				fract.Error((*tks)[oi], "Index accessor is cannot used with non-array variables!")
			}
			val := p.procVal(*vtks)
			i := indexes(v.Val, val, tk)
			var d []obj.Data
			if !v.Val.Arr {
				d = append(d, obj.Data{D: "", T: obj.VStr})
			}
			for _, pos := range i {
				if v.Val.Arr {
					d = append(d, v.Val.D[pos])
				} else {
					if v.Val.D[0].T == obj.VStr {
						d[0].D = d[0].String() + string(v.Val.D[0].String()[pos])
					} else {
						d[0].D = d[0].String() + fmt.Sprint(v.Val.D[0].String()[pos])
					}
				}
			}
			r.Arr = false
			r.D = d
			*r = applyMinus(minus, *r)
			return pos - oi + 1
		case "[":
			// Array initializer.

			// Find close brace.
			ci := pos + 1
			bc := 1
			for ; ci < len(*tks); ci++ {
				tk := (*tks)[ci]
				if tk.T == fract.Brace {
					if tk.Val == "[" {
						bc++
					} else if tk.Val == "]" {
						bc--
						if bc == 0 {
							break
						}
					}
				}
			}
			*r = applyMinus(minus, p.procEnumerableVal((*tks)[pos:ci+1]))
			tks.Rem(pos+1, ci-pos)
			return 0
		case "]":
			// Find open bracket.
			bc := 1
			oi := pos - 1
			for ; oi >= 0; oi-- {
				tk := (*tks)[oi]
				if tk.T == fract.Brace {
					if tk.Val == "]" {
						bc++
					} else if tk.Val == "[" {
						bc--
						if bc == 0 {
							break
						}
					}
				}
			}
			// Finished?
			if oi == 0 {
				r.Arr = true
				r.D = p.procEnumerableVal((*tks)[oi : pos+1]).D
				*r = applyMinus(minus, *r)
				tks.Rem(oi, pos-oi)
				return pos - oi
			}
			endtk := (*tks)[oi-1]
			vi, t, source := p.defByName(endtk)
			if vi == -1 || t != 'v' {
				fract.Error(endtk, "Variable is not defined in this name!: "+endtk.Val)
			}
			vtks := tks.Sub(oi+1, pos-oi-1)
			// Index value is empty?
			if vtks == nil {
				fract.Error(endtk, "Index is not defined!")
			}
			v := source.vars[vi]
			if !v.Val.Arr && v.Val.D[0].T != obj.VStr {
				fract.Error((*tks)[oi], "Index accessor is cannot used with non-array variables!")
			}
			val := p.procVal(*vtks)
			i := indexes(v.Val, val, tk)
			tks.Rem(oi-1, pos-oi+1)
			var d []obj.Data
			if !v.Val.Arr {
				d = append(d, obj.Data{D: "", T: obj.VStr})
			}
			for _, pos := range i {
				if v.Val.Arr {
					d = append(d, v.Val.D[pos])
				} else {
					if v.Val.D[0].T == obj.VStr {
						d[0].D = d[0].String() + string(v.Val.D[0].String()[pos])
					} else {
						d[0].D = d[0].String() + fmt.Sprint(v.Val.D[0].String()[pos])
					}
				}
			}
			r.Arr = len(i) > 1 && d[0].T != obj.VStr || d[0].T == obj.VArray
			r.D = d
			*r = applyMinus(minus, *r)
			return pos - oi + 1
		case ")":
			// Function.

			// Find open parentheses.
			bc := 1
			oi := pos - 1
			for ; oi >= 0; oi-- {
				tk := (*tks)[oi]
				if tk.T == fract.Brace {
					if tk.Val == ")" {
						bc++
					} else if tk.Val == "(" {
						bc--
						if bc == 0 {
							break
						}
					}
				}
			}
			oi--
			v := p.funcCall((*tks)[oi : pos+1])
			if v.D == nil {
				fract.Error((*tks)[oi], "Function is not return any value!")
			}
			*r = applyMinus(minus, v)
			tks.Rem(oi, pos-oi)
			return pos - oi
		}
	}

	//* Single value.
	if strings.HasPrefix(tk.Val, "object.") {
		r.Arr = false
		r.D = []obj.Data{{D: tk.Val, T: obj.VFunc}}
		goto end
	}
	if (tk.T == fract.Value && tk.Val != "true" && tk.Val != "false") && tk.Val[0] != '\'' && tk.Val[0] != '"' {
		if strings.Contains(tk.Val, ".") || strings.ContainsAny(tk.Val, "eE") {
			tk.T = obj.VFloat
		} else {
			tk.T = obj.VInt
		}
		if tk.Val != "NaN" {
			prs, _ := new(big.Float).SetString(tk.Val)
			val, _ := prs.Float64()
			tk.Val = fmt.Sprint(val)
		}
	}
	r.Arr = false
	if tk.Val[0] == '\'' || tk.Val[0] == '"' { // String?
		r.D = []obj.Data{{D: tk.Val[1 : len(tk.Val)-1], T: obj.VStr}}
		tk.T = fract.None // Skip type check.
	} else {
		r.D = []obj.Data{{D: tk.Val}}
	}
	//* Type check.
	if tk.T != fract.None {
		if tk.Val == "true" || tk.Val == "false" {
			r.D[0].T = obj.VBool
			*r = applyMinus(minus, *r)
		} else if tk.T == obj.VFloat { // Float?
			r.D[0].T = obj.VFloat
			*r = applyMinus(minus, *r)
		}
	}
end:
	return 0
}

// Process array value.
func (p *Parser) procArrayVal(tks obj.Tokens) obj.Value {
	v := obj.Value{
		Arr: true,
		D:   []obj.Data{},
	}
	fst := tks[0]
	comma := 1
	bc := 0
	for j := 1; j < len(tks)-1; j++ {
		t := tks[j]
		if t.T == fract.Brace {
			if t.Val == "[" || t.Val == "{" || t.Val == "(" {
				bc++
			} else {
				bc--
			}
		} else if t.T == fract.Comma && bc == 0 {
			lst := tks.Sub(comma, j-comma)
			if lst == nil {
				fract.Error(fst, "Value is not defined!")
			}
			val := p.procVal(*lst)
			if val.Arr {
				v.D = append(v.D, obj.Data{D: val.D, T: obj.VArray})
			} else {
				v.D = append(v.D, val.D...)
			}
			comma = j + 1
		}
	}
	if comma < len(tks)-1 {
		lst := tks.Sub(comma, len(tks)-comma-1)
		if lst == nil {
			fract.Error(fst, "Value is not defined!")
		}
		val := p.procVal(*lst)
		if val.Arr {
			v.D = append(v.D, obj.Data{D: val.D, T: obj.VArray})
		} else {
			v.D = append(v.D, val.D...)
		}
	}
	return v
}

// Process list comprehension.
func (p *Parser) procListComprehension(tks obj.Tokens) obj.Value {
	v := obj.Value{
		Arr: true,
		D:   []obj.Data{},
	}
	var (
		stks obj.Tokens // Select tokens.
		ltks obj.Tokens // Loop tokens.
		ftks obj.Tokens // Filter tokens.
		bc   int
	)
	for i, t := range tks {
		if t.T == fract.Brace {
			if t.Val == "(" || t.Val == "[" || t.Val == "{" {
				bc++
			} else {
				bc--
			}
		}
		if bc > 1 {
			continue
		}
		if t.T == fract.Loop {
			stks = tks[1:i]
		} else if t.T == fract.Comma {
			ltks = tks[len(stks)+1 : i]
			ftks = tks[i+1 : len(tks)-1]
			if len(ftks) == 0 {
				ftks = nil
			}
			break
		}
	}
	if ltks == nil {
		ltks = tks[len(stks)+1 : len(tks)-1]
	}
	if len(ltks) < 2 {
		fract.Error(ltks[0], "Variable name is not defined!")
	}
	nametk := ltks[1]
	// Name is not name?
	if nametk.T != fract.Name {
		fract.Error(nametk, "This is not a valid name!")
	}
	if ln := p.definedName(nametk); ln != -1 {
		fract.Error(nametk, "\""+nametk.Val+"\" is already defined at line: "+fmt.Sprint(ln))
	}
	if len(ltks) < 3 {
		fract.Errorc(ltks[0].F, ltks[0].Ln, ltks[1].Col+len(ltks[1].Val), "Value is not defined!")
	}
	if vtks, inTk := ltks.Sub(3, len(ltks)-3), ltks[2]; vtks != nil {
		ltks = *vtks
	} else {
		fract.Error(inTk, "Value is not defined!")
	}
	varr := p.procVal(ltks)
	// Type is not array?
	if !varr.Arr && varr.D[0].T != obj.VStr {
		fract.Error(ltks[0], "Foreach loop must defined array value!")
	}
	p.vars = append(p.vars, obj.Var{Name: nametk.Val, Val: obj.Value{}})
	vlen := len(p.vars)
	element := &p.vars[vlen-1]
	if element.Name == "_" {
		element.Name = ""
	}
	var length int
	if varr.Arr {
		length = len(varr.D)
	} else {
		length = len(varr.D[0].String())
	}
	// Interpret block.
	for j := 0; j < length; j++ {
		if element.Name != "" {
			if v.Arr {
				element.Val.D = []obj.Data{varr.D[j]}
			} else {
				element.Val.D = []obj.Data{{D: string(varr.D[0].String()[j]), T: obj.VStr}}
			}
		}
		if ftks == nil || p.procCondition(ftks) == "true" {
			val := p.procVal(stks)
			if !val.Arr {
				v.D = append(v.D, val.D...)
			} else {
				v.D = append(v.D, obj.Data{D: val.D, T: obj.VArray})
			}
		}
	}
	p.vars = p.vars[:vlen-1] // Remove variables.
	return v
}

// Process enumerable value.
func (p *Parser) procEnumerableVal(tks obj.Tokens) obj.Value {
	var (
		lc bool
		bc int
	)
	for _, t := range tks {
		if t.T == fract.Brace {
			if t.Val == "(" || t.Val == "[" || t.Val == "{" {
				bc++
			} else {
				bc--
			}
		}
		if bc > 1 {
			continue
		}
		if t.T == fract.Comma {
			break
		} else if !lc && t.T == fract.Loop {
			lc = true
			break
		}
	}
	if lc {
		return p.procListComprehension(tks)
	}
	return p.procArrayVal(tks)
}

// Process value.
func (p *Parser) procVal(tks obj.Tokens) obj.Value {
	p.procRange(&tks)
	// Is conditional expression?
	if j, _ := findConditionOpr(tks); j != -1 {
		return obj.Value{D: []obj.Data{{D: p.procCondition(tks), T: obj.VBool}}}
	}
	v := obj.Value{D: []obj.Data{{}}}
	checkArithmeticProcesses(tks)
	if j := nextopr(tks); j != -1 {
		// Decompose arithmetic operations.
		var opr process
		for j != -1 {
			opr.f = tks[j-1]
			j -= p.procValPart(true, &opr, &tks, j-1)
			opr.opr = tks[j]
			opr.s = tks[j+1]
			j -= p.procValPart(false, &opr, &tks, j+1)
			rv := solveProc(opr)
			opr.opr.Val = "+"
			opr.s = tks[j+1]
			opr.fv = v
			opr.sv = rv
			v = solveProc(opr)
			// Remove processed processes.
			tks.Rem(j-1, 3)
			// Find next operator.
			j = nextopr(tks)
			// Finished?
			if j != -1 && (j == 0 || j == len(tks)-1) {
				opr.fv = v
				opr.opr = tks[j]
				if j == 0 {
					opr.s = tks[j+1]
					p.procValPart(false, &opr, &tks, j+1)
				} else {
					opr.s = tks[j-1]
					p.procValPart(false, &opr, &tks, j-1)
				}
				v = solveProc(opr)
				break
			}
		}
	} else {
		var opr process
		opr.f = tks[0]
		opr.fv.Arr = true //* Ignore nil control if function call.
		p.procValPart(true, &opr, &tks, 0)
		v = opr.fv
	}
	return v
}
