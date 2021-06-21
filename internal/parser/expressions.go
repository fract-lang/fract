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

func compareValues(opr string, d0, d1 obj.Data) bool {
	if d0.T != d1.T && (d0.T == obj.VString || d1.T == obj.VString) {
		return false
	}
	switch opr {
	case "==": // Equals.
		if (d0.T == obj.VString && d0.D != d1.D) ||
			(d0.T != obj.VString && arithmetic.Arithmetic(d0.String()) != arithmetic.Arithmetic(d1.String())) {
			return false
		}
	case "<>": // Not equals.
		if (d0.T == obj.VString && d0.D == d1.D) ||
			(d0.T != obj.VString && arithmetic.Arithmetic(d0.String()) == arithmetic.Arithmetic(d1.String())) {
			return false
		}
	case ">": // Greater.
		if (d0.T == obj.VString && d0.String() <= d1.String()) ||
			(d0.T != obj.VString && arithmetic.Arithmetic(d0.String()) <= arithmetic.Arithmetic(d1.String())) {
			return false
		}
	case "<": // Less.
		if (d0.T == obj.VString && d0.String() >= d1.String()) ||
			(d0.T != obj.VString && arithmetic.Arithmetic(d0.String()) >= arithmetic.Arithmetic(d1.String())) {
			return false
		}
	case ">=": // Greater or equals.
		if (d0.T == obj.VString && d0.String() < d1.String()) ||
			(d0.T != obj.VString && arithmetic.Arithmetic(d0.String()) < arithmetic.Arithmetic(d1.String())) {
			return false
		}
	case "<=": // Less or equals.
		if (d0.T == obj.VString && d0.String() > d1.String()) ||
			(d0.T != obj.VString && arithmetic.Arithmetic(d0.String()) > arithmetic.Arithmetic(d1.String())) {
			return false
		}
	}
	return true
}

func compare(v0, v1 obj.Value, opr obj.Token) bool {
	// In.
	if opr.Val == "in" {
		if !v1.Arr && v1.D[0].T != obj.VString {
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
					if d.T != obj.VString {
						fract.Error(opr, "All datas is not string!")
					}
					if strings.Contains(dt, d.String()) {
						return true
					}
				}
			} else {
				if v1.D[0].T != obj.VString {
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
		d0 := v0.D[0]
		d1 := v1.D[0]
		if (d0.T == obj.VString && d1.T != obj.VString) || (d0.T != obj.VString && d1.T == obj.VString) {
			fract.Error(opr, "The in keyword should use with string or enumerable data types!")
		}
		return compareValues(opr.Val, d0, d1)
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
			if !compareValues(opr.Val, d, v1.D[i]) {
				return false
			}
		}
		return true
	}
	// Single value comparison.
	return compareValues(opr.Val, v0.D[0], v1.D[0])
}

// processCondition returns condition result.
func (p *Parser) processCondition(tks obj.Tokens) string {
	p.processRange(&tks)
	T := obj.Value{D: []obj.Data{{D: "true"}}}
	// Process condition.
	ors := decomposeConditionalProcess(tks, "||")
	for _, or := range *ors {
		// Decompose and conditions.
		ands := decomposeConditionalProcess(or, "&&")
		// Is and long statement?
		if len(*ands) > 1 {
			for _, and := range *ands {
				opri, opr := findConditionOperator(and)
				// Operator is not found?
				if opri == -1 {
					opr.Val = "=="
					if compare(p.processValue(and), T, opr) {
						return "true"
					}
					continue
				}
				// Operator is first or last?
				if opri == 0 {
					fract.Error(and[0], "Comparison values are missing!")
				} else if opri == len(and)-1 {
					fract.Error(and[len(and)-1], "Comparison values are missing!")
				}
				if !compare(
					p.processValue(*and.Sub(0, opri)), p.processValue(*and.Sub(opri+1, len(and)-opri-1)), opr) {
					return "false"
				}
			}
			return "true"
		}
		opri, opr := findConditionOperator(or)
		// Operator is not found?
		if opri == -1 {
			opr.Val = "=="
			if compare(p.processValue(or), T, opr) {
				return "true"
			}
			continue
		}
		// Operator is first or last?
		if opri == 0 {
			fract.Error(or[0], "Comparison values are missing!")
		} else if opri == len(or)-1 {
			fract.Error(or[len(or)-1], "Comparison values are missing!")
		}
		if compare(p.processValue(*or.Sub(0, opri)), p.processValue(*or.Sub(opri+1, len(or)-opri-1)), opr) {
			return "true"
		}
	}
	return "false"
}

// Get string arithmetic compatible data.
func arith(tks obj.Token, d obj.Data) string {
	ret := d.String()
	switch d.T {
	case obj.VFunction:
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

// processRange by value processor principles.
func (p *Parser) processRange(tks *obj.Tokens) {
	for {
		rg, pos := decomposeBrace(tks, "(", ")", true)
		/* Parentheses are not found! */
		if pos == -1 {
			return
		}
		val := p.processValue(rg)
		if val.Arr {
			tks.Insert(pos, obj.Token{Val: "[", T: fract.Brace})
			for _, current := range val.D {
				pos++
				tks.Insert(pos, obj.Token{Val: current.Format(), T: fract.Value})
				pos++
				tks.Insert(pos, obj.Token{Val: ",", T: fract.Comma})
			}
			pos++
			tks.Insert(pos, obj.Token{Val: "]", T: fract.Brace})
		} else {
			if val.D[0].T == obj.VString {
				tks.Insert(pos, obj.Token{Val: "\"" + val.D[0].String() + "\"", T: fract.Value})
			} else {
				tks.Insert(pos, obj.Token{
					Val: val.D[0].Format(),
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
func solve(opr obj.Token, f, s float64) float64 {
	var r float64
	if opr.Val == "\\" || opr.Val == "\\\\" { // Divide with bigger.
		if opr.Val == "\\" {
			opr.Val = "/"
		} else {
			opr.Val = "//"
		}
		if f < s {
			cache := f
			f = s
			s = cache
		}
	}
	switch opr.Val {
	case "+": // Addition.
		r = f + s
	case "-": // Subtraction.
		r = f - s
	case "*": // Multiply.
		r = f * s
	case "/", "//": // Division.
		if f == 0 || s == 0 {
			fract.Error(opr, "Divide by zero!")
		}
		r = f / s
		if opr.Val == "//" {
			r = math.RoundToEven(r)
		}
	case "|": // Binary or.
		r = float64(int(f) | int(s))
	case "&": // Binary and.
		r = float64(int(f) & int(s))
	case "^": // Bitwise exclusive or.
		r = float64(int(f) ^ int(s))
	case "**": // Exponentiation.
		r = math.Pow(f, s)
	case "%": // Mod.
		r = math.Mod(f, s)
	case "<<": // Left shift.
		if s < 0 {
			fract.Error(opr, "Shifter is cannot should be negative!")
		}
		r = float64(int(f) << int(s))
	case ">>": // Right shift.
		if s < 0 {
			fract.Error(opr, "Shifter is cannot should be negative!")
		}
		r = float64(int(f) >> int(s))
	default:
		fract.Error(opr, "Operator is invalid!")
	}
	return r
}

// Check data and set ready.
func readyData(p process, d obj.Data) obj.Data {
	if p.fv.D[0].T == obj.VString || p.sv.D[0].T == obj.VString {
		d.T = obj.VString
	} else if p.opr.Val == "/" || p.opr.Val == "\\" ||
		p.fv.D[0].T == obj.VFloat || p.sv.D[0].T == obj.VFloat {
		d.T = obj.VFloat
		d.D = d.Format()
		return d
	}
	return d
}

// solveProcess solve arithmetic process.
func solveProcess(p process) obj.Value {
	v := obj.Value{D: []obj.Data{{}}}
	// String?
	if (len(p.fv.D) != 0 && p.fv.D[0].T == obj.VString) || (len(p.sv.D) != 0 && p.sv.D[0].T == obj.VString) {
		if p.fv.D[0].T == p.sv.D[0].T { // Both string?
			v.D[0].T = obj.VString
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

		v.D[0].T = obj.VString
		if p.fv.D[0].T == obj.VString {
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
				if p.sv.D[0].T != obj.VInteger {
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
				if p.fv.D[0].T != obj.VInteger {
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
			ar := arithmetic.Arithmetic(arith(p.opr, f.D[0]))
			for i, d := range s.D {
				if d.T == obj.VArray {
					s.D[i] = readyData(p, obj.Data{
						D: solveProcess(process{
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
					s.D[i] = readyData(p, obj.Data{D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, ar, arithmetic.Arithmetic(arith(p.opr, d))))})
				}
			}
			v.D = s.D
		} else {
			for i, f := range p.fv.D {
				s := p.sv.D[i]
				if f.T == obj.VArray || s.T == obj.VArray {
					proc := process{
						f:   p.f,
						s:   p.s,
						opr: p.opr,
					}
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
					p.fv.D[i] = readyData(p, obj.Data{D: solveProcess(proc).D, T: obj.VArray})
				} else {
					p.fv.D[i] = readyData(p,
						obj.Data{
							D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, arithmetic.Arithmetic(arith(p.opr, f)), arithmetic.Arithmetic(s.String()))),
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
		ar := arithmetic.Arithmetic(arith(p.opr, s.D[0]))
		for i, d := range f.D {
			if d.T == obj.VArray {
				f.D[i] = readyData(p, obj.Data{
					D: solveProcess(process{
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
						D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, arithmetic.Arithmetic(arith(p.opr, d)), ar)),
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
				D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, arithmetic.Arithmetic(arith(p.opr, p.fv.D[0])), arithmetic.Arithmetic(arith(p.opr, p.sv.D[0])))),
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
			if d.T == obj.VBoolean || d.T == obj.VFloat || d.T == obj.VInteger {
				d.D = fmt.Sprintf(fract.FloatFormat, -arithmetic.Arithmetic(d.String()))
				val.D[i].D = d.Format()
			}
		}
		return val
	}
	if d := val.D[0]; d.T == obj.VBoolean || d.T == obj.VFloat || d.T == obj.VInteger {
		d.D = fmt.Sprintf(fract.FloatFormat, -arithmetic.Arithmetic(d.String()))
		val.D[0].D = d.Format()
	}
	return val
}

func (p *Parser) processOperationValue(fst bool, opr *process, tks *obj.Tokens, pos int) int {
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
					vi, t, src := p.defineByName(tk)
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
					val := p.processValue(*vtks)
					if val.Arr {
						fract.Error((*tks)[pos], "Arrays is not used in index access!")
					} else if val.D[0].T != obj.VInteger {
						fract.Error((*tks)[pos], "Only integer values can used in index access!")
					}
					vp, err := strconv.Atoi(arith((*vtks)[0], val.D[0]))
					if err != nil {
						fract.Error((*tks)[pos], "Invalid value!")
					}
					v := src.vars[vi]
					if !v.Val.Arr && v.Val.D[0].T != obj.VString {
						fract.Error((*tks)[pos], "Index accessor is cannot used with non-array variables!")
					}
					if v.Val.Arr {
						vp = processIndex(len(v.Val.D), vp)
					} else {
						vp = processIndex(len(v.Val.D[0].String()), vp)
					}
					if vp == -1 {
						fract.Error((*tks)[pos], "Index is out of range!")
					}
					tks.Remove(pos+1, ci-pos)
					var d obj.Data
					if v.Val.Arr {
						d = v.Val.D[vp]
					} else {
						if v.Val.D[0].T == obj.VString {
							d = obj.Data{D: string(v.Val.D[0].String()[vp]), T: obj.VString}
						} else {
							d = obj.Data{D: fmt.Sprint(v.Val.D[0].String()[vp])}
						}
					}
					r.Arr = d.T == obj.VArray
					if r.Arr {
						r.D = d.D.([]obj.Data)
					} else {
						r.D = []obj.Data{d}
					}
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
					v := p.processFunctionCall((*tks)[pos:ci])
					if !opr.fv.Arr && v.D == nil {
						fract.Error(tk, "Function is not return any value!")
					}
					tks.Remove(pos+1, ci-pos-1)
					*r = applyMinus(minus, v)
					return 0
				}
			}
		}
		vi, t, src := p.defineByName(tk)
		if vi == -1 {
			fract.Error(tk, "Variable is not defined in this name: "+tk.Val)
		}
		switch t {
		case 'f':
			*r = obj.Value{
				D: []obj.Data{{D: src.funcs[vi], T: obj.VFunction}},
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
				r.D = p.processArrayValue(*tks.Sub(oi, pos-oi+1)).D
				*r = applyMinus(minus, *r)
				tks.Remove(oi, pos-oi)
				return pos - oi
			}
			endtk := (*tks)[oi-1]
			vi, t, src := p.defineByName(tk)
			if vi == -1 || t != 'v' {
				fract.Error(endtk, "Variable is not defined in this name: "+endtk.Val)
			}
			vtks := tks.Sub(oi+1, pos-oi-1)
			// Index value is empty?
			if vtks == nil {
				fract.Error(endtk, "Index is not defined!")
			}
			val := p.processValue(*vtks)
			if val.Arr {
				fract.Error((*tks)[pos], "Arrays is not used in index access!")
			} else if val.D[0].T != obj.VInteger {
				fract.Error((*tks)[pos], "Only integer values can used in index access!")
			}
			vp, err := strconv.Atoi(arith((*vtks)[0], val.D[0]))
			if err != nil {
				fract.Error((*tks)[oi], "Invalid value!")
			}
			v := src.vars[vi]
			if !v.Val.Arr && v.Val.D[0].T != obj.VString {
				fract.Error((*tks)[oi], "Index accessor is cannot used with non-array variables!")
			}
			if v.Val.Arr {
				vp = processIndex(len(v.Val.D), vp)
			} else {
				vp = processIndex(len(v.Val.D[0].String()), vp)
			}
			if vp == -1 {
				fract.Error((*tks)[oi], "Index is out of range!")
			}
			tks.Remove(oi-1, pos-oi+1)
			var d obj.Data
			if v.Val.Arr {
				d = v.Val.D[vp]
			} else {
				if v.Val.D[0].T == obj.VString {
					d = obj.Data{D: string(v.Val.D[0].String()[vp]), T: obj.VString}
				} else {
					d = obj.Data{D: fmt.Sprint(v.Val.D[0].String()[vp])}
				}
			}
			r.D = []obj.Data{d}
			r.Arr = false
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
			*r = applyMinus(minus, p.processArrayValue((*tks)[pos:ci+1]))
			tks.Remove(pos+1, ci-pos)
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
				r.D = p.processArrayValue((*tks)[oi : pos+1]).D
				*r = applyMinus(minus, *r)
				tks.Remove(oi, pos-oi)
				return pos - oi
			}
			endtk := (*tks)[oi-1]
			vi, t, source := p.defineByName(endtk)
			if vi == -1 || t != 'v' {
				fract.Error(endtk, "Variable is not defined in this name!: "+endtk.Val)
			}
			vtks := tks.Sub(oi+1, pos-oi-1)
			// Index value is empty?
			if vtks == nil {
				fract.Error(endtk, "Index is not defined!")
			}
			val := p.processValue(*vtks)
			if val.Arr {
				fract.Error((*tks)[pos], "Arrays is not used in index access!")
			} else if val.D[0].T != obj.VInteger {
				fract.Error((*tks)[pos], "Only integer values can used in index access!")
			}
			vp, err := strconv.Atoi(arith(tk, val.D[0]))
			if err != nil {
				fract.Error((*tks)[oi], "Invalid value!")
			}
			v := source.vars[vi]
			if !v.Val.Arr && v.Val.D[0].T != obj.VString {
				fract.Error((*tks)[oi], "Index accessor is cannot used with non-array variables!")
			}
			if v.Val.Arr {
				vp = processIndex(len(v.Val.D), vp)
			} else {
				vp = processIndex(len(v.Val.D[0].String()), vp)
			}
			if vp == -1 {
				fract.Error((*tks)[oi], "Index is out of range!")
			}
			tks.Remove(oi-1, pos-oi+1)
			var d obj.Data
			if v.Val.Arr {
				d = v.Val.D[vp]
			} else {
				if v.Val.D[0].T == obj.VString {
					d = obj.Data{D: string(v.Val.D[0].String()[vp]), T: obj.VString}
				} else {
					d = obj.Data{D: fmt.Sprint(v.Val.D[0].String()[vp])}
				}
			}
			r.Arr = d.T == obj.VArray
			if r.Arr {
				r.D = d.D.([]obj.Data)
			} else {
				r.D = []obj.Data{d}
			}
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
			v := p.processFunctionCall((*tks)[oi : pos+1])
			if v.D == nil {
				fract.Error((*tks)[oi], "Function is not return any value!")
			}
			*r = applyMinus(minus, v)
			tks.Remove(oi, pos-oi)
			return pos - oi
		}
	}

	//* Single value.
	if strings.HasPrefix(tk.Val, "object.") {
		fract.Error(tk, "\""+tk.Val+"\" is not compatible with arithmetic processes!")
	}
	if (tk.T == fract.Value && tk.Val != "true" && tk.Val != "false") && tk.Val[0] != '\'' && tk.Val[0] != '"' {
		if strings.Contains(tk.Val, ".") || strings.ContainsAny(tk.Val, "eE") {
			tk.T = obj.VFloat
		} else {
			tk.T = obj.VInteger
		}
		if tk.Val != "NaN" {
			prs, _ := new(big.Float).SetString(tk.Val)
			val, _ := prs.Float64()
			tk.Val = fmt.Sprint(val)
		}
	}
	r.Arr = false
	if tk.Val[0] == '\'' || tk.Val[0] == '"' { // String?
		r.D = []obj.Data{{D: tk.Val[1 : len(tk.Val)-1], T: obj.VString}}
		tk.T = fract.None // Skip type check.
	} else {
		r.D = []obj.Data{{D: tk.Val}}
	}
	//* Type check.
	if tk.T != fract.None {
		if tk.Val == "true" || tk.Val == "false" {
			r.D[0].T = obj.VBoolean
			*r = applyMinus(minus, *r)
		} else if tk.T == obj.VFloat { // Float?
			r.D[0].T = obj.VFloat
			*r = applyMinus(minus, *r)
		}
	}
	return 0
}

func (p *Parser) processArrayValue(tks obj.Tokens) obj.Value {
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
			val := p.processValue(*lst)
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
		val := p.processValue(*lst)
		if val.Arr {
			v.D = append(v.D, obj.Data{D: val.D, T: obj.VArray})
		} else {
			v.D = append(v.D, val.D...)
		}
	}
	return v
}

func (p *Parser) processValue(tks obj.Tokens) obj.Value {
	p.processRange(&tks)
	v := obj.Value{D: []obj.Data{{}}}
	// Is conditional expression?
	if j, _ := findConditionOperator(tks); j != -1 {
		v.D = []obj.Data{{D: p.processCondition(tks), T: obj.VBoolean}}
		return v
	}
	checkArithmeticProcesses(tks)
	if j := indexProcess(tks); j != -1 {
		// Decompose arithmetic operations.
		var opr process
		for j != -1 {
			opr.f = tks[j-1]
			j -= p.processOperationValue(true, &opr, &tks, j-1)
			opr.opr = tks[j]
			opr.s = tks[j+1]
			j -= p.processOperationValue(false, &opr, &tks, j+1)
			resultValue := solveProcess(opr)
			opr.opr.Val = "+"
			opr.s = tks[j+1]
			opr.fv = v
			opr.sv = resultValue
			v = solveProcess(opr)
			// Remove processed processes.
			tks.Remove(j-1, 3)
			tks.Insert(j-1, obj.Token{Val: "0"})
			// Find next operator.
			j = indexProcess(tks)
		}
	} else {
		var opr process
		opr.f = tks[0]
		opr.fv.Arr = true //* Ignore nil control if function call.
		p.processOperationValue(true, &opr, &tks, 0)
		v = opr.fv
	}
	return v
}