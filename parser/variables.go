package parser

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
	"github.com/fract-lang/fract/pkg/value"
)

// Metadata of variable declaration.
type varinfo struct {
	sdec      bool
	constant  bool
	mut       bool
	protected bool
}

// Append variable to source.
func (p *Parser) varadd(md varinfo, tks obj.Tokens) {
	name := tks[0]
	if strings.Contains(name.V, ".") {
		fract.IPanic(name, obj.SyntaxPanic, "Names is cannot include dot!")
	} else if name.V == "_" {
		fract.IPanic(name, obj.SyntaxPanic, "Ignore operator is cannot be variable name!")
	}
	// Name is already defined?
	if ln := p.definedName(name); ln != -1 {
		fract.IPanic(name, obj.NamePanic, "\""+name.V+"\" is already defined at line: "+fmt.Sprint(ln))
	}
	tksLen := len(tks)
	// Setter is not defined?
	if tksLen < 2 {
		fract.IPanicC(name.F, name.Ln, name.Col+len(name.V), obj.SyntaxPanic, "Setter is not found!")
	}
	setter := tks[1]
	// Setter is not a setter operator?
	if setter.T != fract.Operator || (setter.V != "=" && !md.sdec || setter.V != ":=" && md.sdec) {
		fract.IPanic(setter, obj.SyntaxPanic, "Invalid setter operator: "+setter.V)
	}
	// Value is not defined?
	if tksLen < 3 {
		fract.IPanicC(setter.F, setter.Ln, setter.Col+len(setter.V), obj.SyntaxPanic, "Value is not given!")
	}
	v := p.procVal(*tks.Sub(2, tksLen-2))
	if v.D == nil {
		fract.IPanic(tks[2], obj.ValuePanic, "Invalid value!")
	}
	if p.funcTempVars != -1 {
		p.funcTempVars++
	}
	p.vars = append(p.vars,
		obj.Var{
			Name:      name.V,
			V:         v,
			Ln:        name.Ln,
			Const:     md.constant,
			Mut:       md.mut,
			Protected: md.protected,
		})
}

// Process variable declaration.
func (p *Parser) vardec(tks obj.Tokens, protected bool) {
	// Name is not defined?
	if len(tks) < 2 {
		first := tks[0]
		fract.IPanicC(first.F, first.Ln, first.Col+len(first.V), obj.SyntaxPanic, "Name is not given!")
	}
	md := varinfo{
		constant:  tks[0].V == "const",
		mut:       tks[0].V == "mut",
		protected: protected,
	}
	pre := tks[1]
	if pre.T == fract.Name {
		p.varadd(md, tks[1:])
	} else if pre.T == fract.Brace && pre.V == "(" {
		tks = tks[2 : len(tks)-1]
		lst := 0
		ln := tks[0].Ln
		bc := 0
		for j, t := range tks {
			if t.T == fract.Brace {
				switch t.V {
				case "{", "[", "(":
					bc++
				default:
					bc--
					ln = t.Ln
				}
			}
			if bc > 0 {
				continue
			}
			if ln < t.Ln {
				p.varadd(md, tks[lst:j])
				lst = j
				ln = t.Ln
			}
		}
		if len(tks) != lst {
			p.varadd(md, tks[lst:])
		}
	} else {
		fract.IPanic(pre, obj.SyntaxPanic, "Invalid syntax!")
	}
}

// Process short variable declaration.
func (p *Parser) varsdec(tks obj.Tokens) {
	// Name is not defined?
	if len(tks) < 2 {
		first := tks[0]
		fract.IPanicC(first.F, first.Ln, first.Col+len(first.V), obj.SyntaxPanic, "Name is not given!")
	}
	if tks[0].T != fract.Name {
		fract.IPanic(tks[0], obj.SyntaxPanic, "Invalid syntax!")
	}
	var md varinfo
	md.sdec = true
	p.varadd(md, tks)
}

// Process variable set statement.
func (p *Parser) varset(tks obj.Tokens) {
	name := tks[0]
	// Name is not name?
	if name.T != fract.Name {
		fract.IPanic(name, obj.SyntaxPanic, "Invalid name!")
	} else if name.V == "_" {
		fract.IPanic(name, obj.SyntaxPanic, "Ignore operator is cannot set!")
	}
	j, _ := p.varIndexByName(name)
	if j == -1 {
		fract.IPanic(name, obj.NamePanic, "Variable is not defined in this name: "+name.V)
	}
	v := p.vars[j]
	// Check const state.
	if v.Const {
		fract.IPanic(tks[1], obj.SyntaxPanic, "Values is cannot changed of constant defines!")
	}
	var (
		setpos []int
		setter = tks[1]
	)
	// TODO: Add Map.
	// Array setter?
	if setter.T == fract.Brace && setter.V == "[" {
		// Variable is not array?
		if v.V.T != value.Array && v.V.T != value.Map && v.V.T != value.Str {
			fract.IPanic(setter, obj.ValuePanic, "Variable is not enumerable!")
		}
		bc := 1
		// Find close bracket.
		for j := 2; j < len(tks); j++ {
			t := tks[j]
			if t.T == fract.Brace {
				switch t.V {
				case "[":
					bc++
				case "]":
					bc--
				}
			}
			if bc > 0 {
				continue
			}
			vtks := tks.Sub(2, j-2)
			// Index value is empty?
			if vtks == nil {
				fract.IPanic(setter, obj.SyntaxPanic, "Index is not given!")
			}
			tks = append(obj.Tokens{tks[1]}, tks[j+1:]...)
			setpos = selections(v.V, p.procVal(*vtks), tks[0]).([]int)
			break
		}
	}
	setter = tks[1]
	// Value are not defined?
	if len(tks) < 3 {
		fract.IPanicC(setter.F, setter.Ln, setter.Col+len(setter.V), obj.SyntaxPanic, "Value is not given!")
	}
	val := p.procVal(*tks.Sub(2, len(tks)-2))
	if val.D == nil {
		fract.IPanic(tks[2], obj.ValuePanic, "Invalid value!")
	}
	if setpos != nil {
		for _, pos := range setpos {
			switch setter.V {
			case "=": // =
				if v.V.T == value.Array {
					v.V.D.([]value.Val)[pos] = val
					break
				}
				if val.T != value.Str {
					fract.IPanic(setter, obj.ValuePanic, "Value type is not string!")
				} else if len(val.String()) > 1 {
					fract.IPanic(setter, obj.ValuePanic, "Value length is should be maximum one!")
				}
				bytes := []byte(v.V.String())
				if val.D == "" {
					bytes[pos] = 0
				} else {
					bytes[pos] = val.String()[0]
				}
				v.V.D = string(bytes)
			default: // Other assignments.
				if v.V.T == value.Array {
					v.V.D.([]value.Val)[pos] = solveProc(process{
						opr: obj.Token{V: string(setter.V[:len(setter.V)-1])},
						f:   tks,
						fv:  v.V.D.([]value.Val)[pos],
						s:   obj.Tokens{setter},
						sv:  val,
					})
					break
				}
				val = solveProc(process{
					opr: obj.Token{V: string(setter.V[:len(setter.V)-1])},
					f:   tks,
					fv:  value.Val{D: v.V.D.(string)[pos], T: value.Int},
					s:   obj.Tokens{setter},
					sv:  val,
				})
				if val.T != value.Str {
					fract.IPanic(setter, obj.ValuePanic, "Value type is not string!")
				} else if len(val.String()) > 1 {
					fract.IPanic(setter, obj.ValuePanic, "Value length is should be maximum one!")
				}
				bytes := []byte(v.V.String())
				if val.D == "" {
					bytes[pos] = 0
				} else {
					bytes[pos] = val.String()[0]
				}
				v.V.D = string(bytes)
			}
		}
	} else {
		switch setter.V {
		case "=": // =
			v.V = val
		default: // Other assignments.
			v.V = solveProc(process{
				opr: obj.Token{V: string(setter.V[:len(setter.V)-1])},
				f:   tks,
				fv:  v.V,
				s:   obj.Tokens{setter},
				sv:  val,
			})
		}
	}
	p.vars[j] = v
}
