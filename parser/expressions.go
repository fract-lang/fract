package parser

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
	"github.com/fract-lang/fract/pkg/value"
)

// Compare arithmetic values.
func compVals(opr string, d0, d1 value.Data) bool {
	if d0.T != d1.T && (d0.T == value.Str || d1.T == value.Str) {
		return false
	}
	switch opr {
	case "==": // Equals.
		if !d0.Equals(d1) {
			return false
		}
	case "<>": // Not equals.
		if !d0.NotEquals(d1) {
			return false
		}
	case ">": // Greater.
		if !d0.Greater(d1) {
			return false
		}
	case "<": // Less.
		if !d0.Less(d1) {
			return false
		}
	case ">=": // Greater or equals.
		if !d0.GreaterEquals(d1) {
			return false
		}
	case "<=": // Less or equals.
		if !d0.LessEquals(d1) {
			return false
		}
	}
	return true
}

// Compare values.
func comp(v0, v1 value.Val, opr obj.Token) bool {
	// In.
	if opr.V == "in" {
		if !v1.Arr && v1.D[0].T != value.Str {
			fract.IPanic(opr, obj.ValuePanic, "Value is not enumerable!")
		}
		if v1.Arr {
			dt := v0.String()
			for _, d := range v1.D {
				if strings.Contains(d.String(), dt) {
					return true
				}
			}
			return false
		}
		// String.
		if v0.Arr {
			dt := v1.D[0].String()
			for _, d := range v0.D {
				if d.T != value.Str {
					fract.IPanic(opr, obj.ValuePanic, "All values is not string!")
				}
				if strings.Contains(dt, d.String()) {
					return true
				}
			}
		} else {
			if v1.D[0].T != value.Str {
				fract.IPanic(opr, obj.ValuePanic, "All datas is not string!")
			}
			if strings.Contains(v1.D[0].String(), v0.D[0].String()) {
				return true
			}
		}
		return false
	}
	// Array comparison.
	if v0.Arr || v1.Arr {
		if (v0.Arr && !v1.Arr) || (!v0.Arr && v1.Arr) {
			return false
		}
		if len(v0.D) != len(v1.D) {
			return opr.V == "<>"
		}
		for i, d := range v0.D {
			if !compVals(opr.V, d, v1.D[i]) {
				return false
			}
		}
		return true
	}
	// Single value comparison.
	d0, d1 := v0.D[0], v1.D[0]
	if (d0.T == value.Str && d1.T != value.Str) || (d0.T != value.Str && d1.T == value.Str) {
		fract.IPanic(opr, obj.ValuePanic, "The in keyword should use with string or enumerable data types!")
	}
	return compVals(opr.V, d0, d1)
}

// procCondition returns condition result.
func (p *Parser) procCondition(tks obj.Tokens) string {
	T := value.Val{D: []value.Data{{D: "true"}}}
	// Process condition.
	ors := conditionalProcesses(tks, "||")
	for _, or := range ors {
		// Decompose and conditions.
		ands := conditionalProcesses(or, "&&")
		// Is and long statement?
		if len(ands) > 1 {
			for _, and := range ands {
				i, opr := findConditionOpr(and)
				// Operator is not found?
				if i == -1 {
					opr.V = "=="
					if comp(p.procVal(and), T, opr) {
						return "true"
					}
					return "false"
				}
				// Operator is first or last?
				if i == 0 {
					fract.IPanic(and[0], obj.SyntaxPanic, "Comparison values are missing!")
				} else if i == len(and)-1 {
					fract.IPanic(and[len(and)-1], obj.SyntaxPanic, "Comparison values are missing!")
				}
				if !comp(p.procVal(and[:i]), p.procVal(*and.Sub(i+1, len(and)-i-1)), opr) {
					return "false"
				}
			}
			return "true"
		}
		i, opr := findConditionOpr(or)
		// Operator is not found?
		if i == -1 {
			opr.V = "=="
			if comp(p.procVal(or), T, opr) {
				return "true"
			}
			continue
		}
		// Operator is first or last?
		if i == 0 {
			fract.IPanic(or[0], obj.SyntaxPanic, "Comparison values are missing!")
		} else if i == len(or)-1 {
			fract.IPanic(or[len(or)-1], obj.SyntaxPanic, "Comparison values are missing!")
		}
		if comp(p.procVal(or[:i]), p.procVal(*or.Sub(i+1, len(or)-i-1)), opr) {
			return "true"
		}
	}
	return "false"
}

// Get string arithmetic compatible data.
func arith(tks obj.Token, d value.Data) string {
	ret := d.String()
	switch d.T {
	case value.Func:
		fract.IPanic(tks, obj.ArithmeticPanic, "\""+ret+"\" is not compatible with arithmetic processes!")
	}
	return ret
}

// process instance for solver.
type process struct {
	f   obj.Tokens // Tokens of first value.
	fv  value.Val  // Value instance of first value.
	s   obj.Tokens // Tokens of second value.
	sv  value.Val  // Value instance of second value.
	opr obj.Token  // Operator of process.
}

// solve process.
func solve(opr obj.Token, a, b float64) float64 {
	var r float64
	switch opr.V {
	case "+": // Addition.
		r = a + b
	case "-": // Subtraction.
		r = a - b
	case "*": // Multiply.
		r = a * b
	case "/", "//": // Division.
		if a == 0 || b == 0 {
			fract.Panic(opr, obj.DivideByZeroPanic, "Divide by zero!")
		}
		r = a / b
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
			fract.IPanic(opr, obj.ArithmeticPanic, "Shifter is cannot should be negative!")
		}
		r = float64(int(a) << int(b))
	case ">>": // Right shift.
		if b < 0 {
			fract.IPanic(opr, obj.ArithmeticPanic, "Shifter is cannot should be negative!")
		}
		r = float64(int(a) >> int(b))
	default:
		fract.IPanic(opr, obj.SyntaxPanic, "Operator is invalid!")
	}
	return r
}

// Check data and set ready.
func readyData(p process, d value.Data) value.Data {
	if p.fv.D[0].T == value.Str || p.sv.D[0].T == value.Str {
		d.T = value.Str
	} else if p.opr.V == "/" || p.fv.D[0].T == value.Float || p.sv.D[0].T == value.Float {
		d.T = value.Float
		return d
	}
	return d
}

// solveProc solve arithmetic process.
func solveProc(p process) value.Val {
	v := value.Val{D: []value.Data{{D: "0", T: value.Int}}}
	// String?
	if (len(p.fv.D) != 0 && p.fv.D[0].T == value.Str) || (len(p.sv.D) != 0 && p.sv.D[0].T == value.Str) {
		if p.fv.D[0].T == p.sv.D[0].T { // Both string?
			v.D[0].T = value.Str
			switch p.opr.V {
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
				fract.IPanic(p.opr, obj.ArithmeticPanic, "This operator is not defined for string types!")
			}
			return v
		}

		v.D[0].T = value.Str
		if p.sv.D[0].T == value.Str {
			p.fv, p.sv = p.sv, p.fv
		}
		if p.sv.Arr {
			if len(p.sv.D) == 0 {
				v.D = p.fv.D
				return v
			}
			if len(p.fv.D[0].String()) != len(p.sv.D) && (len(p.fv.D[0].String()) != 1 && len(p.sv.D) != 1) {
				fract.IPanic(p.s[0], obj.ArithmeticPanic, "Array element count is not one or equals to first array!")
			}
			if strings.Contains(p.sv.D[0].String(), ".") {
				fract.IPanic(p.s[0], obj.ArithmeticPanic, "Only string and integer values can concatenate string values!")
			}
			r, _ := strconv.ParseInt(p.sv.D[0].String(), 10, 64)
			rn := rune(r)
			var sb strings.Builder
			for _, r := range p.fv.D[0].String() {
				switch p.opr.V {
				case "+":
					sb.WriteByte(byte(r + rn))
				case "-":
					sb.WriteByte(byte(r - rn))
				default:
					fract.IPanic(p.opr, obj.ArithmeticPanic, "This operator is not defined for string types!")
				}
			}
			v.D[0].D = sb.String()
		} else {
			if p.sv.D[0].T != value.Int {
				fract.IPanic(p.s[0], obj.ArithmeticPanic, "Only string and integer values can concatenate string values!")
			}
			var s string
			rs, _ := strconv.ParseInt(p.sv.D[0].String(), 10, 64)
			rn := byte(rs)
			for _, r := range p.fv.D[0].String() {
				switch p.opr.V {
				case "+":
					s += string(byte(r) + rn)
				case "-":
					s += string(byte(r) - rn)
				default:
					fract.IPanic(p.opr, obj.ArithmeticPanic, "This operator is not defined for string types!")
				}
			}
			v.D[0].D = s
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
			fract.IPanic(p.s[0], obj.ArithmeticPanic, "Array element count is not one or equals to first array!")
		}
		if len(p.fv.D) == 1 || len(p.sv.D) == 1 {
			f, s := p.fv, p.sv
			if len(f.D) != 1 {
				f, s = s, f
			}
			ar := value.Conv(arith(p.opr, f.D[0]))
			for i, d := range s.D {
				if d.T == value.Array {
					s.D[i] = readyData(p, value.Data{
						D: solveProc(process{
							f:   p.f,
							fv:  s,
							s:   p.s,
							sv:  value.Val{D: d.D.([]value.Data), Arr: true},
							opr: p.opr,
						}).D,
						T: value.Array,
					})
				} else {
					s.D[i] = readyData(p, value.Data{D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, ar, value.Conv(arith(p.opr, d))))})
				}
			}
			v.D = s.D
		} else {
			for i, f := range p.fv.D {
				s := p.sv.D[i]
				if f.T == value.Array || s.T == value.Array {
					proc := process{f: p.f, s: p.s, opr: p.opr}
					if f.T == value.Array {
						proc.fv = value.Val{D: f.D.([]value.Data), Arr: true}
					} else {
						proc.fv = value.Val{D: []value.Data{f}}
					}
					if s.T == value.Array {
						proc.sv = value.Val{D: s.D.([]value.Data), Arr: true}
					} else {
						proc.sv = value.Val{D: []value.Data{s}}
					}
					p.fv.D[i] = readyData(p, value.Data{D: solveProc(proc).D, T: value.Array})
				} else {
					p.fv.D[i] = readyData(p,
						value.Data{
							D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, value.Conv(arith(p.opr, f)), value.Conv(s.String()))),
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
		ar := value.Conv(arith(p.opr, s.D[0]))
		for i, d := range f.D {
			if d.T == value.Array {
				f.D[i] = readyData(p, value.Data{
					D: solveProc(process{
						f:   p.f,
						fv:  s,
						s:   p.s,
						sv:  value.Val{D: d.D.([]value.Data), Arr: true},
						opr: p.opr,
					}).D,
					T: value.Array,
				})
			} else {
				f.D[i] = readyData(p,
					value.Data{
						D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, value.Conv(arith(p.opr, d)), ar)),
					})
			}
		}
		v.D = f.D
	} else {
		if len(p.fv.D) == 0 {
			p.fv.D = []value.Data{{D: "0", T: value.Int}}
		}
		v.D[0] = readyData(p,
			value.Data{
				D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, value.Conv(arith(p.opr, p.fv.D[0])), value.Conv(arith(p.opr, p.sv.D[0])))),
			})
	}
	return v
}

// applyMinus operator.
func applyMinus(minus obj.Token, v value.Val) value.Val {
	if minus.V[0] != '-' {
		return v
	}
	val := value.Val{Arr: v.Arr, D: append([]value.Data{}, v.D...)}
	for i, d := range val.D {
		switch d.T {
		case value.Bool, value.Float, value.Int:
			val.D[i].D = fmt.Sprintf(fract.FloatFormat, -value.Conv(d.String()))
		default:
			fract.IPanic(minus, obj.ArithmeticPanic, "Bad operand type for unary!")
		}
	}
	return val
}

func (p *Parser) selectArrayElems(v value.Val, indexes []int) value.Val {
	var r value.Val
	if !v.Arr {
		r.D = append(r.D, value.Data{D: "", T: value.Str})
	}
	if len(indexes) == 1 {
		d := v.D[indexes[0]]
		if d.T == value.Array {
			r.D = d.D.([]value.Data)
			r.Arr = true
		} else {
			r.D = []value.Data{d}
		}
	} else {
		for _, pos := range indexes {
			if v.Arr {
				r.D = append(r.D, v.D[pos])
			} else {
				if v.D[0].T == value.Str {
					r.D[0].D = r.D[0].String() + string(v.D[0].String()[pos])
				} else {
					r.D[0].D = r.D[0].String() + fmt.Sprint(v.D[0].String()[pos])
				}
			}
		}
		r.Arr = len(indexes) > 1 && r.D[0].T != value.Str || r.D[0].T == value.Array
	}
	return r
}

// Process value part.
func (p *Parser) procValPart(nilch bool, tks obj.Tokens) value.Val {
	var (
		r  = value.Val{}
		tk = tks[0]
	)
	// Single value.
	if len(tks) == 1 {
		if tk.T == fract.Name {
			vi, t, src := p.defByName(tk)
			if vi == -1 {
				fract.IPanic(tk, obj.NamePanic, "Name is not defined: "+tk.V)
			}
			switch t {
			case 'f': // Function.
				r = value.Val{D: []value.Data{{D: src.funcs[vi], T: value.Func}}}
			case 'v': // Value.
				v := src.vars[vi]
				val := v.V
				if !v.Mut { //! Immutability.
					val.D = append(make([]value.Data, 0), v.V.D...)
				}
				r = applyMinus(tk, val)
			}
		} else if tk.V[0] == '\'' || tk.V[0] == '"' { // String?
			r.D = []value.Data{{D: tk.V[1 : len(tk.V)-1], T: value.Str}}
		} else if tk.T == fract.Value && tk.V == "true" || tk.V == "false" {
			r.D = []value.Data{{D: tk.V, T: value.Bool}}
		} else if tk.T == fract.Value {
			if strings.Contains(tk.V, ".") || strings.ContainsAny(tk.V, "eE") {
				tk.T = value.Float
			} else {
				tk.T = value.Int
			}
			if tk.V != "NaN" {
				prs, _ := new(big.Float).SetString(tk.V)
				val, _ := prs.Float64()
				tk.V = fmt.Sprint(val)
			}
			r.D = []value.Data{{D: tk.V, T: tk.T}}
		} else if strings.HasPrefix(tk.V, "object.") {
			r.D = []value.Data{{D: tk.V, T: value.Func}}
		} else {
			fract.IPanic(tk, obj.ValuePanic, "Invalid value!")
		}
		return r
	}
	switch i, tk := len(tks)-1, tks[len(tks)-1]; tk.T {
	case fract.Brace:
		bc := 0
		switch tk.V {
		case ")":
			var vtks obj.Tokens
			for ; i >= 0; i-- {
				t := tks[i]
				if t.T != fract.Brace {
					continue
				}
				switch t.V {
				case ")":
					bc++
				case "(":
					bc--
				}
				if bc > 0 {
					continue
				}
				vtks = tks[:i]
				break
			}
			if len(vtks) == 0 && bc == 0 {
				tk, tks = tks[0], tks[1:len(tks)-1]
				if len(tks) == 0 {
					fract.IPanic(tk, obj.SyntaxPanic, "Invalid syntax!")
				}
				return applyMinus(tk, p.procVal(tks))
			}
			// Function call.
			v := p.procValPart(nilch, vtks)
			if v.Arr || v.D[0].T != value.Func {
				fract.IPanic(tks[len(vtks)], obj.ValuePanic, "Value is not function!")
			}
			return applyMinus(tk, p.funcCallModel(v.D[0].D.(function), tks[len(vtks):]).call())
		case "]":
			var vtks obj.Tokens
			for ; i >= 0; i-- {
				t := tks[i]
				if t.T != fract.Brace {
					continue
				}
				switch t.V {
				case "]":
					bc++
				case "[":
					bc--
				}
				if bc > 0 {
					continue
				}
				vtks = tks[:i]
				break
			}
			if len(vtks) == 0 && bc == 0 {
				return applyMinus(tk, p.procEnumerableVal(tks))
			}
			v := p.procValPart(nilch, vtks)
			if !v.Arr && v.D[0].T != value.Str {
				fract.IPanic(tk, obj.ValuePanic, "Index accessor is cannot used with non-array variables!")
			}
			return applyMinus(tk, p.selectArrayElems(v, indexes(v, p.procVal(tks[len(vtks):]), tk)))
		case "}":
			var vtks obj.Tokens
			for ; i >= 0; i-- {
				t := tks[i]
				if t.T != fract.Brace {
					continue
				}
				switch t.V {
				case "}":
					bc++
				case "{":
					bc--
				}
				if bc > 0 {
					continue
				}
				vtks = tks[:i]
				break
			}
			l := len(vtks)
			if l == 0 && bc == 0 {
				fract.IPanic(tk, obj.SyntaxPanic, "Invalid syntax!")
			} else if vtks[0].T != fract.Func {
				fract.IPanic(tk, obj.SyntaxPanic, "Invalid syntax!")
			} else if l > 1 && (vtks[1].T != fract.Brace || vtks[1].V != "(") {
				fract.IPanic(vtks[1], obj.SyntaxPanic, "Invalid syntax!")
			} else if l > 1 && (vtks[l-1].T != fract.Brace || vtks[l-1].V != ")") {
				fract.IPanic(vtks[l-1], obj.SyntaxPanic, "Invalid syntax!")
			}
			f := function{
				name: "anonymous",
				src:  p,
				tks:  p.getBlock(tks[len(vtks):]),
			}
			if l > 1 {
				vtks = vtks[1:]
				p.setFuncParams(&f, &vtks)
			}
			r.D = []value.Data{{D: f, T: value.Func}}
			return r
		}
	}
	fract.IPanic(tk, obj.ValuePanic, "Invalid value!")
	return r
}

// Process array value.
func (p *Parser) procArrayVal(tks obj.Tokens) value.Val {
	v := value.Val{
		Arr: true,
		D:   []value.Data{},
	}
	fst := tks[0]
	comma := 1
	bc := 0
	for j := 1; j < len(tks)-1; j++ {
		switch t := tks[j]; t.T {
		case fract.Brace:
			switch t.V {
			case "{", "[", "(":
				bc++
			default:
				bc--
			}
		case fract.Comma:
			if bc != 0 {
				break
			}
			lst := tks.Sub(comma, j-comma)
			if lst == nil {
				fract.IPanic(fst, obj.SyntaxPanic, "Value is not given!")
			}
			val := p.procVal(*lst)
			if val.Arr {
				v.D = append(v.D, value.Data{D: val.D, T: value.Array})
			} else {
				v.D = append(v.D, val.D...)
			}
			comma = j + 1
		}
	}
	if comma < len(tks)-1 {
		lst := tks.Sub(comma, len(tks)-comma-1)
		if lst == nil {
			fract.IPanic(fst, obj.SyntaxPanic, "Value is not given!")
		}
		val := p.procVal(*lst)
		if val.Arr {
			v.D = append(v.D, value.Data{D: val.D, T: value.Array})
		} else {
			v.D = append(v.D, val.D...)
		}
	}
	return v
}

// Process list comprehension.
func (p *Parser) procListComprehension(tks obj.Tokens) value.Val {
	v := value.Val{Arr: true, D: []value.Data{}}
	var (
		stks obj.Tokens // Select tokens.
		ltks obj.Tokens // Loop tokens.
		ftks obj.Tokens // Filter tokens.
		bc   int
	)
	for i, t := range tks {
		if t.T == fract.Brace {
			switch t.V {
			case "{", "[", "(":
				bc++
			default:
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
		fract.IPanic(ltks[0], obj.SyntaxPanic, "Variable name is not given!")
	}
	nametk := ltks[1]
	// Name is not name?
	if nametk.T != fract.Name {
		fract.IPanic(nametk, obj.SyntaxPanic, "This is not a valid name!")
	}
	if ln := p.definedName(nametk); ln != -1 {
		fract.IPanic(nametk, obj.NamePanic, "\""+nametk.V+"\" is already defined at line: "+fmt.Sprint(ln))
	}
	if len(ltks) < 3 {
		fract.IPanicC(ltks[0].F, ltks[0].Ln, ltks[1].Col+len(ltks[1].V), obj.SyntaxPanic, "Value is not given!")
	}
	if vtks, inTk := ltks.Sub(3, len(ltks)-3), ltks[2]; vtks != nil {
		ltks = *vtks
	} else {
		fract.IPanic(inTk, obj.SyntaxPanic, "Value is not given!")
	}
	varr := p.procVal(ltks)
	// Type is not array?
	if !varr.Arr && varr.D[0].T != value.Str {
		fract.IPanic(ltks[0], obj.ValuePanic, "Foreach loop must defined array value!")
	}
	p.vars = append(p.vars, obj.Var{Name: nametk.V, V: value.Val{}})
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
				element.V.D = []value.Data{varr.D[j]}
			} else {
				element.V.D = []value.Data{{D: string(varr.D[0].String()[j]), T: value.Str}}
			}
		}
		if ftks == nil || p.procCondition(ftks) == "true" {
			val := p.procVal(stks)
			if !val.Arr {
				v.D = append(v.D, val.D...)
			} else {
				v.D = append(v.D, value.Data{D: val.D, T: value.Array})
			}
		}
	}
	p.vars = p.vars[:vlen-1] // Remove variables.
	return v
}

// Process enumerable value.
func (p *Parser) procEnumerableVal(tks obj.Tokens) value.Val {
	var (
		lc bool
		bc int
	)
	for _, t := range tks {
		if t.T == fract.Brace {
			switch t.V {
			case "{", "[", "(":
				bc++
			default:
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
func (p *Parser) procVal(tks obj.Tokens) value.Val {
	// Is conditional expression?
	if j, _ := findConditionOpr(tks); j != -1 {
		return value.Val{D: []value.Data{{D: p.procCondition(tks), T: value.Bool}}}
	}
	procs := arithmeticProcesses(tks)
	if len(procs) == 1 {
		return p.procValPart(false, procs[0])
	}
	var v value.Val
	var opr process
	j := nextopr(procs)
	for j != -1 {
		opr.f = procs[j-1]
		opr.fv = p.procValPart(true, opr.f)
		opr.opr = procs[j][0]
		opr.s = procs[j+1]
		opr.sv = p.procValPart(true, opr.s)
		rv := solveProc(opr)
		if v.D != nil {
			opr.opr.V = "+"
			opr.s = procs[j+1]
			opr.fv = v
			opr.sv = rv
			v = solveProc(opr)
		} else {
			v = rv
		}
		// Remove computed processes.
		procs = append(procs[:j-1], procs[j+2:]...)
		// Find next operator.
		j = nextopr(procs)
		// If last value to compute.
		if j != -1 && (j == 0 || j == len(procs)-1) {
			opr.fv = v
			opr.opr = procs[j][0]
			if j == 0 {
				opr.s = procs[j+1]
			} else {
				opr.s = procs[j-1]
			}
			opr.fv = p.procValPart(true, opr.s)
			v = solveProc(opr)
			break
		}
	}
	return v
}
