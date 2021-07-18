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
func compVals(opr string, d0, d1 value.Val) bool {
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
		if v1.T != value.Array && v1.T != value.Str {
			fract.IPanic(opr, obj.ValuePanic, "Value is can should be string or array!")
		}
		if v1.T == value.Array {
			dt := v0.String()
			for _, d := range v1.D.([]value.Val) {
				if strings.Contains(d.String(), dt) {
					return true
				}
			}
			return false
		}
		// String.
		if v0.T == value.Array {
			dt := v1.String()
			for _, d := range v0.D.([]value.Val) {
				if d.T != value.Str {
					fract.IPanic(opr, obj.ValuePanic, "All values is not string!")
				}
				if strings.Contains(dt, d.String()) {
					return true
				}
			}
		} else {
			if v1.T != value.Str {
				fract.IPanic(opr, obj.ValuePanic, "All datas is not string!")
			}
			if strings.Contains(v1.String(), v0.String()) {
				return true
			}
		}
		return false
	}
	return compVals(opr.V, v0, v1)
}

// procCondition returns condition result.
func (p *Parser) procCondition(tks obj.Tokens) string {
	T := value.Val{D: "true", T: value.Bool}
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
func arith(tks obj.Token, d value.Val) string {
	ret := d.String()
	switch d.T {
	case value.Func:
		fract.IPanic(tks, obj.ArithmeticPanic, "\""+ret+"\" is not compatible with arithmetic processes!")
	case value.Map:
		fract.IPanic(tks, obj.ArithmeticPanic, "\"object.map\" is not compatible with arithmetic processes!")
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
func readyData(p process, d value.Val) value.Val {
	if p.fv.T == value.Str || p.sv.T == value.Str {
		d.T = value.Str
	} else if p.opr.V == "/" || p.fv.T == value.Float || p.sv.T == value.Float {
		d.T = value.Float
		return d
	}
	return d
}

// solveProc solve arithmetic process.
func solveProc(p process) value.Val {
	v := value.Val{D: "0", T: value.Int}
	fl := p.fv.Len()
	sl := p.sv.Len()
	// String?
	if (fl != 0 && p.fv.T == value.Str) || (sl != 0 && p.sv.T == value.Str) {
		if p.fv.T == p.sv.T { // Both string?
			v.T = value.Str
			switch p.opr.V {
			case "+":
				v.D = p.fv.String() + p.sv.String()
			case "-":
				flen := len(p.fv.String())
				slen := len(p.sv.String())
				if flen == 0 || slen == 0 {
					v.D = ""
					break
				}
				if flen == 1 && slen > 1 {
					r, _ := strconv.ParseInt(p.fv.String(), 10, 32)
					fr := rune(r)
					for _, r := range p.sv.String() {
						v.D = v.String() + string(fr-r)
					}
				} else if slen == 1 && flen > 1 {
					r, _ := strconv.ParseInt(p.sv.String(), 10, 32)
					fr := rune(r)
					for _, r := range p.fv.String() {
						v.D = v.String() + string(fr-r)
					}
				} else {
					for i, r := range p.fv.String() {
						v.D = v.String() + string(r-rune(p.sv.String()[i]))
					}
				}
			default:
				fract.IPanic(p.opr, obj.ArithmeticPanic, "This operator is not defined for string types!")
			}
			return v
		}

		v.T = value.Str
		if p.sv.T == value.Str {
			p.fv, p.sv = p.sv, p.fv
		}
		if p.sv.T == value.Array {
			if sl == 0 {
				v.D = p.fv.D
				return v
			}
			if len(p.fv.String()) != sl && (len(p.fv.String()) != 1 && sl != 1) {
				fract.IPanic(p.s[0], obj.ArithmeticPanic, "Array element count is not one or equals to first array!")
			}
			if strings.Contains(p.sv.String(), ".") {
				fract.IPanic(p.s[0], obj.ArithmeticPanic, "Only string and integer values can concatenate string values!")
			}
			r, _ := strconv.ParseInt(p.sv.String(), 10, 64)
			rn := rune(r)
			var sb strings.Builder
			for _, r := range p.fv.String() {
				switch p.opr.V {
				case "+":
					sb.WriteByte(byte(r + rn))
				case "-":
					sb.WriteByte(byte(r - rn))
				default:
					fract.IPanic(p.opr, obj.ArithmeticPanic, "This operator is not defined for string types!")
				}
			}
			v.D = sb.String()
		} else {
			if p.sv.T != value.Int {
				fract.IPanic(p.s[0], obj.ArithmeticPanic, "Only string and integer values can concatenate string values!")
			}
			var s string
			rs, _ := strconv.ParseInt(p.sv.String(), 10, 64)
			rn := byte(rs)
			for _, r := range p.fv.String() {
				switch p.opr.V {
				case "+":
					s += string(byte(r) + rn)
				case "-":
					s += string(byte(r) - rn)
				default:
					fract.IPanic(p.opr, obj.ArithmeticPanic, "This operator is not defined for string types!")
				}
			}
			v.D = s
		}
		return v
	}

	if p.fv.T == value.Array && p.sv.T == value.Array {
		v.T = value.Array
		if fl == 0 {
			v.D = p.sv.D
			return v
		} else if sl == 0 {
			v.D = p.fv.D
			return v
		}
		if fl != sl && fl != 1 && sl != 1 {
			fract.IPanic(p.s[0], obj.ArithmeticPanic, "Array element count is not one or equals to first array!")
		}
		if fl == 1 || sl == 1 {
			f, s := p.fv, p.sv
			if f.Len() != 1 {
				f, s = s, f
			}
			ar := value.Conv(arith(p.opr, f.D.([]value.Val)[0]))
			for i, d := range s.D.([]value.Val) {
				if d.T == value.Array {
					s.D.([]value.Val)[i] = readyData(p, value.Val{
						D: solveProc(process{
							f:   p.f,
							fv:  s,
							s:   p.s,
							sv:  d,
							opr: p.opr,
						}).D,
						T: value.Array,
					})
				} else {
					s.D.([]value.Val)[i] = readyData(p, value.Val{
						D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, ar, value.Conv(arith(p.opr, d)))),
					})
				}
			}
			v.D = s.D
		} else {
			for i, f := range p.fv.D.([]value.Val) {
				s := p.sv.D.([]value.Val)[i]
				if f.T == value.Array || s.T == value.Array {
					proc := process{f: p.f, s: p.s, opr: p.opr}
					if f.T == value.Array {
						proc.fv = value.Val{D: f.D, T: value.Array}
					} else {
						proc.fv = value.Val{D: f}
					}
					if s.T == value.Array {
						proc.sv = value.Val{D: s.D, T: value.Array}
					} else {
						proc.sv = value.Val{D: s}
					}
					p.fv.D.([]value.Val)[i] = readyData(p, value.Val{D: solveProc(proc).D, T: value.Array})
				} else {
					p.fv.D.([]value.Val)[i] = readyData(p, value.Val{
						D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, value.Conv(arith(p.opr, f)), value.Conv(s.String()))),
					})
				}
			}
			v.D = p.fv.D
		}
	} else if p.fv.T == value.Array || p.sv.T == value.Array {
		v.T = value.Array
		if p.fv.T == value.Array && fl == 0 {
			v.D = p.sv.D
			return v
		} else if p.sv.T == value.Array && sl == 0 {
			v.D = p.fv.D
			return v
		}
		f, s := p.fv, p.sv
		if f.T != value.Array {
			f, s = s, f
		}
		ar := value.Conv(arith(p.opr, s))
		for i, d := range f.D.([]value.Val) {
			if d.T == value.Array {
				f.D.([]value.Val)[i] = readyData(p, solveProc(process{
					f:   p.f,
					fv:  s,
					s:   p.s,
					sv:  d,
					opr: p.opr,
				}))
			} else {
				f.D.([]value.Val)[i] = readyData(p, value.Val{
					D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, value.Conv(arith(p.opr, d)), ar)),
				})
			}
		}
		v = f
	} else {
		v = readyData(p,
			value.Val{
				D: fmt.Sprintf(fract.FloatFormat, solve(p.opr, value.Conv(arith(p.opr, p.fv)), value.Conv(arith(p.opr, p.sv)))),
			})
	}
	return v
}

// applyMinus operator.
func applyMinus(minus obj.Token, v value.Val) value.Val {
	if minus.V[0] != '-' {
		return v
	}
	for i, d := range v.D.([]value.Val) {
		switch d.T {
		case value.Bool, value.Float, value.Int:
			v.D.([]value.Val)[i].D = fmt.Sprintf(fract.FloatFormat, -value.Conv(d.String()))
		default:
			fract.IPanic(minus, obj.ArithmeticPanic, "Bad operand type for unary!")
		}
	}
	return v
}

// Select enumerable object elements.
func (p *Parser) selectEnum(v value.Val, tk obj.Token, s interface{}) value.Val {
	var r value.Val
	switch v.T {
	case value.Array:
		i := s.([]int)
		if len(i) == 1 {
			return v.D.([]value.Val)[i[0]]
		}
		r = value.Val{D: []value.Val{}, T: value.Array}
		for _, pos := range i {
			r.D = append(r.D.([]value.Val), v.D.([]value.Val)[pos])
		}
	case value.Map:
		m := v.D.(value.MapModel)
		switch t := s.(type) {
		case []value.Val:
			rm := value.MapModel{}
			for _, k := range t {
				d, ok := m[k]
				if !ok {
					fract.IPanic(tk, obj.ValuePanic, "Key is not exists!")
				}
				rm[k] = d
			}
			r = value.Val{D: rm, T: value.Map}
		case value.Val:
			d, ok := m[t]
			if !ok {
				fract.IPanic(tk, obj.ValuePanic, "Key is not exists!")
			}
			return d
		}
	case value.Str:
		r = value.Val{D: "", T: value.Str}
		for _, pos := range s.([]int) {
			r.D = r.String() + string(v.String()[pos])
		}
	}
	return r
}

type valPartInfo struct {
	tks obj.Tokens
	mut bool // Force to mutability.
}

// Process value part.
func (p *Parser) procValPart(i valPartInfo) value.Val {
	if i.tks[0].T == fract.Var && i.tks[0].V == "mut" {
		if len(i.tks) == 1 {
			fract.IPanic(i.tks[0], obj.SyntaxPanic, "Value is not given!")
		}
		i.mut = true
		i.tks = i.tks[1:]
		return p.procValPart(i)
	}
	var (
		tk = i.tks[0]
	)
	// Single value.
	if len(i.tks) == 1 {
		if tk.T == fract.Name {
			vi, t, src := p.defByName(tk)
			if vi == -1 {
				fract.IPanic(tk, obj.NamePanic, "Name is not defined: "+tk.V)
			}
			switch t {
			case 'f': // Function.
				return value.Val{D: src.funcs[vi], T: value.Func}
			case 'v': // Value.
				v := src.vars[vi]
				var val value.Val
				if !v.Mut && !i.mut { //! Immutability.
					val = v.V.Immut()
				} else {
					val = v.V
				}
				return applyMinus(tk, val)
			}
		} else if tk.V[0] == '\'' || tk.V[0] == '"' {
			return value.Val{D: tk.V[1 : len(tk.V)-1], T: value.Str}
		} else if tk.V == "true" || tk.V == "false" {
			return value.Val{D: tk.V, T: value.Bool}
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
			return value.Val{D: tk.V, T: tk.T}
		} else if strings.HasPrefix(tk.V, "object.func") {
			return value.Val{D: tk.V, T: value.Func}
		} else {
			fract.IPanic(tk, obj.ValuePanic, "Invalid value!")
		}
	}
	switch j, tk := len(i.tks)-1, i.tks[len(i.tks)-1]; tk.T {
	case fract.Brace:
		bc := 0
		switch tk.V {
		case ")":
			var vtks obj.Tokens
			for ; j >= 0; j-- {
				t := i.tks[j]
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
				vtks = i.tks[:j]
				break
			}
			if len(vtks) == 0 && bc == 0 {
				tk, i.tks = i.tks[0], i.tks[1:len(i.tks)-1]
				if len(i.tks) == 0 {
					fract.IPanic(tk, obj.SyntaxPanic, "Invalid syntax!")
				}
				return applyMinus(tk, p.procVal(i.tks))
			}
			// Function call.
			v := p.procValPart(valPartInfo{tks: vtks})
			if v.T != value.Func {
				fract.IPanic(i.tks[len(vtks)], obj.ValuePanic, "Value is not function!")
			}
			return applyMinus(tk, p.funcCallModel(v.D.(function), i.tks[len(vtks):]).call())
		case "]":
			var vtks obj.Tokens
			for ; j >= 0; j-- {
				t := i.tks[j]
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
				vtks = i.tks[:j]
				break
			}
			if len(vtks) == 0 && bc == 0 {
				return applyMinus(tk, p.procEnumerableVal(i.tks))
			}
			v := p.procValPart(valPartInfo{tks: vtks})
			if v.T != value.Array && v.T != value.Map && v.T != value.Str {
				fract.IPanic(vtks[0], obj.ValuePanic, "Index accessor is cannot used with not enumerable values!")
			}
			return applyMinus(tk, p.selectEnum(v, tk, selections(v, p.procVal(i.tks[len(vtks)+1:len(i.tks)-1]), tk)))
		case "}":
			var vtks obj.Tokens
			for ; j >= 0; j-- {
				t := i.tks[j]
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
				vtks = i.tks[:j]
				break
			}
			l := len(vtks)
			if l == 0 && bc == 0 || vtks[0].T != fract.Func {
				return applyMinus(tk, p.procEnumerableVal(i.tks))
			} else if l > 1 && (vtks[1].T != fract.Brace || vtks[1].V != "(") {
				fract.IPanic(vtks[1], obj.SyntaxPanic, "Invalid syntax!")
			} else if l > 1 && (vtks[l-1].T != fract.Brace || vtks[l-1].V != ")") {
				fract.IPanic(vtks[l-1], obj.SyntaxPanic, "Invalid syntax!")
			}
			f := function{
				name: "anonymous",
				src:  p,
				tks:  p.getBlock(i.tks[len(vtks):]),
			}
			if f.tks == nil {
				f.tks = []obj.Tokens{}
			}
			if l > 1 {
				vtks = vtks[1:]
				p.setFuncParams(&f, &vtks)
			}
			return value.Val{D: f, T: value.Func}
		}
	}
	fract.IPanic(tk, obj.ValuePanic, "Invalid value!")
	return value.Val{}
}

// Process array value.
func (p *Parser) procArrayVal(tks obj.Tokens) value.Val {
	v := value.Val{D: []value.Val{}, T: value.Array}
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
			v.D = append(v.D.([]value.Val), val)
			comma = j + 1
		}
	}
	if comma < len(tks)-1 {
		lst := tks.Sub(comma, len(tks)-comma-1)
		if lst == nil {
			fract.IPanic(fst, obj.SyntaxPanic, "Value is not given!")
		}
		val := p.procVal(*lst)
		v.D = append(v.D.([]value.Val), val)
	}
	return v
}

// Process map value.
func (p *Parser) procMapVal(tks obj.Tokens) value.Val {
	fst := tks[0]
	comma := 1
	bc := 0
	m := value.MapModel{}
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
			var (
				i  int
				l  int = len(*lst)
				tk obj.Token
			)
			for i, tk = range *lst {
				switch tk.T {
				case fract.Brace:
					switch tk.V {
					case "{", "[", "(":
						bc++
					default:
						bc--
					}
				case fract.Colon:
					if bc != 0 {
						break
					}
					if i+1 >= l {
						fract.IPanic(tk, obj.SyntaxPanic, "Value is not given!")
					}
					key := p.procVal((*lst)[:i])
					if key.T == value.Array {
						_, ok := m[key]
						if ok {
							fract.IPanic(tk, obj.ValuePanic, "Key is already defined!")
						}
						m[key] = p.procVal((*lst)[i+1:])
					} else {
						_, ok := m[key]
						if ok {
							fract.IPanic(tk, obj.ValuePanic, "Key is already defined!")
						}
						m[key] = p.procVal((*lst)[i+1:])
					}
					comma = j + 1
					lst = nil
				}
			}
			if lst != nil {
				fract.IPanic((*lst)[l-1], obj.SyntaxPanic, "Value identifier is not found!")
			}
		}
	}
	if comma < len(tks)-1 {
		lst := *tks.Sub(comma, len(tks)-comma-1)
		i := -1
		l := len(lst)
		for j, tk := range lst {
			switch tk.T {
			case fract.Brace:
				switch tk.V {
				case "{", "[", "(":
					bc++
				default:
					bc--
				}
			case fract.Colon:
				if bc != 0 {
					break
				}
				i = j
			}
			if i != -1 {
				break
			}
		}
		if i+1 >= l {
			fract.IPanic(lst[i], obj.SyntaxPanic, "Value is not given!")
		}
		key := p.procVal(lst[:i])
		if key.T == value.Array {
			_, ok := m[key]
			if ok {
				fract.IPanic(lst[i], obj.ValuePanic, "Key is already defined!")
			}
			m[key] = p.procVal(lst[i+1:])
		} else {
			_, ok := m[key]
			if ok {
				fract.IPanic(lst[i], obj.ValuePanic, "Key is already defined!")
			}
			m[key] = p.procVal(lst[i+1:])
		}
		lst = nil
	}
	v := value.Val{D: m, T: value.Map}
	return v
}

// Process list comprehension.
func (p *Parser) procListComprehension(tks obj.Tokens) value.Val {
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
	if !varr.IsEnum() {
		fract.IPanic(ltks[0], obj.ValuePanic, "Foreach loop must defined enumerable value!")
	}
	p.vars = append(p.vars, obj.Var{Name: nametk.V})
	vlen := len(p.vars)
	element := &p.vars[vlen-1]
	if element.Name == "_" {
		element.Name = ""
	}
	// Interpret block.
	v := value.Val{D: []value.Val{}, T: value.Array}
	l := loop{enum: varr}
	l.run(func() {
		element.V = l.b
		if ftks == nil || p.procCondition(ftks) == "true" {
			val := p.procVal(stks)
			v.D = append(v.D.([]value.Val), val)
		}
	})
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
			if tks[0].V != "[" {
				fract.IPanic(tks[0], obj.SyntaxPanic, "Invalid syntax!")
			}
			lc = true
			break
		}
	}
	if lc {
		return p.procListComprehension(tks)
	} else if tks[0].V == "{" {
		return p.procMapVal(tks)
	}
	return p.procArrayVal(tks)
}

// Process value.
func (p *Parser) procVal(tks obj.Tokens) value.Val {
	// Is conditional expression?
	if j, _ := findConditionOpr(tks); j != -1 {
		return value.Val{D: p.procCondition(tks), T: value.Bool}
	}
	procs := arithmeticProcesses(tks)
	if len(procs) == 1 {
		return p.procValPart(valPartInfo{tks: procs[0]})
	}
	var v value.Val
	var opr process
	j := nextopr(procs)
	for j != -1 {
		opr.f = procs[j-1]
		opr.fv = p.procValPart(valPartInfo{tks: opr.f})
		if opr.fv.T == fract.None {
			fract.IPanic(opr.f[0], obj.ValuePanic, "Value is not given!")
		}
		opr.opr = procs[j][0]
		opr.s = procs[j+1]
		opr.sv = p.procValPart(valPartInfo{tks: opr.s})
		if opr.sv.T == fract.None {
			fract.IPanic(opr.s[0], obj.ValuePanic, "Value is not given!")
		}
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
			opr.fv = p.procValPart(valPartInfo{tks: opr.s})
			v = solveProc(opr)
			break
		}
	}
	return v
}
