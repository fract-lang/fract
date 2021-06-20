package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// Metadata of variable declaration.
type varinfo struct {
	constant  bool
	mut       bool
	protected bool
}

// appendVariable to source from tokens.
func (p *Parser) appendVariable(md varinfo, tks obj.Tokens) {
	name := tks[0]
	if strings.Contains(name.Val, ".") {
		fract.Error(name, "Names is cannot include dot!")
	} else if name.Val == "_" {
		fract.Error(name, "Ignore operator is cannot be variable name!")
	}
	// Name is already defined?
	if ln := p.definedName(name); ln != -1 {
		fract.Error(name, "\""+name.Val+"\" is already defined at line: "+fmt.Sprint(ln))
	}
	tksLen := len(tks)
	// Setter is not defined?
	if tksLen < 2 {
		fract.Errorc(name.F, name.Ln, name.Col+len(name.Val), "Setter is not found!")
	}
	setter := tks[1]
	// Setter is not a setter operator?
	if setter.T != fract.Operator && setter.Val != "=" {
		fract.Error(setter, "This is not a setter operator: "+setter.Val)
	}
	// Value is not defined?
	if tksLen < 3 {
		fract.Errorc(setter.F, setter.Ln, setter.Col+len(setter.Val), "Value is not defined!")
	}
	v := p.processValue(*tks.Sub(2, tksLen-2))
	if v.D == nil {
		fract.Error(tks[2], "Invalid value!")
	}
	p.vars = append(p.vars,
		obj.Var{
			Name:      name.Val,
			Val:       v,
			Ln:        name.Ln,
			Const:     md.constant,
			Mut:       md.mut,
			Protected: md.protected,
		})
}

func (p *Parser) processVariableDeclaration(tks []obj.Token, protected bool) {
	// Name is not defined?
	if len(tks) < 2 {
		first := tks[0]
		fract.Errorc(first.F, first.Ln, first.Col+len(first.Val), "Name is not found!")
	}
	md := varinfo{
		constant:  tks[0].Val == "const",
		mut:       tks[0].Val == "mut",
		protected: protected,
	}
	pre := tks[1]
	if pre.T == fract.Name {
		p.appendVariable(md, tks[1:])
	} else if pre.T == fract.Brace && pre.Val == "(" {
		tks = tks[2 : len(tks)-1]
		lst := 0
		ln := tks[0].Ln
		bc := 0
		for j, t := range tks {
			if t.T == fract.Brace {
				if t.Val == "{" || t.Val == "[" || t.Val == "(" {
					bc++
				} else {
					bc--
					ln = t.Ln
				}
			}
			if bc > 0 {
				continue
			}
			if ln < t.Ln {
				p.appendVariable(md, tks[lst:j])
				lst = j
				ln = t.Ln
			}
		}
		if len(tks) != lst {
			p.appendVariable(md, tks[lst:])
		}
	} else {
		fract.Error(pre, "Invalid syntax!")
	}
}

// Process variable set statement.
func (p *Parser) processVariableSet(tks obj.Tokens) {
	name := tks[0]
	// Name is not name?
	if name.T != fract.Name {
		fract.Error(name, "This is not a valid name!")
	} else if name.Val == "_" {
		fract.Error(name, "Ignore operator is cannot set!")
	}
	j, _ := p.variableIndexByName(name)
	if j == -1 {
		fract.Error(name, "Variable is not defined in this name: "+name.Val)
	}
	v := p.vars[j]
	// Check const state.
	if v.Const {
		fract.Error(tks[1], "Values is cannot changed of constant defines!")
	}
	setter := tks[1]
	setpos := -1
	// Array setter?
	if setter.T == fract.Brace && setter.Val == "[" {
		// Variable is not array?
		if !v.Val.Arr && v.Val.D[0].T != obj.VString {
			fract.Error(setter, "Variable is not array!")
		}
		// Find close bracket.
		for j := 2; j < len(tks); j++ {
			t := tks[j]
			if t.T != fract.Brace || t.Val != "]" {
				continue
			}
			vtks := tks.Sub(2, j-2)
			// Index value is empty?
			if vtks == nil {
				fract.Error(setter, "Index is not defined!")
			}
			pos, err := strconv.Atoi(p.processValue(*vtks).D[0].String())
			if err != nil {
				fract.Error(setter, "Value out of range!")
			}
			if v.Val.Arr {
				pos = processIndex(len(v.Val.D), pos)
			} else {
				pos = processIndex(len(v.Val.D[0].String()), pos)
			}
			if pos == -1 {
				fract.Error(setter, "Index is out of range!")
			}
			setpos = pos
			tks.Remove(1, j)
			setter = tks[1]
			break
		}
	}
	// Value are not defined?
	if len(tks) < 3 {
		fract.Errorc(setter.F, setter.Ln, setter.Col+len(setter.Val), "Value is not defined!")
	}
	val := p.processValue(*tks.Sub(2, len(tks)-2))
	if val.D == nil {
		fract.Error(tks[2], "Invalid value!")
	}
	if setpos != -1 {
		if val.Arr {
			fract.Error(setter, "Array is cannot set as indexed value!")
		}
		switch setter.Val {
		case "=": // =
			if v.Val.Arr {
				v.Val.D[setpos] = val.D[0]
			} else {
				if val.D[0].T != obj.VString {
					fract.Error(setter, "Value type is not string!")
				} else if len(val.D[0].String()) > 1 {
					fract.Error(setter, "Value length is should be maximum one!")
				}
				bytes := []byte(v.Val.D[0].String())
				if val.D[0].D == "" {
					bytes[setpos] = 0
				} else {
					bytes[setpos] = val.D[0].String()[0]
				}
				v.Val.D[0].D = string(bytes)
			}
		default: // Other assignments.
			if v.Val.Arr {
				v.Val.D[setpos] = solveProcess(
					process{
						opr: obj.Token{Val: string(setter.Val[:len(setter.Val)-1])},
						f:   tks[0],
						fv:  obj.Value{D: []obj.Data{v.Val.D[setpos]}},
						s:   setter,
						sv:  val,
					}).D[0]
			} else {
				val = solveProcess(
					process{
						opr: obj.Token{Val: string(setter.Val[:len(setter.Val)-1])},
						f:   tks[0],
						fv:  obj.Value{D: []obj.Data{v.Val.D[setpos]}},
						s:   setter,
						sv:  val,
					})
				if val.D[0].T != obj.VString {
					fract.Error(setter, "Value type is not string!")
				} else if len(val.D[0].String()) > 1 {
					fract.Error(setter, "Value length is should be maximum one!")
				}
				bytes := []byte(v.Val.D[0].String())
				if val.D[0].D == "" {
					bytes[setpos] = 0
				} else {
					bytes[setpos] = val.D[0].String()[0]
				}
				v.Val.D[0].D = string(bytes)
			}
		}
	} else {
		switch setter.Val {
		case "=": // =
			v.Val = val
		default: // Other assignments.
			v.Val = solveProcess(
				process{
					opr: obj.Token{Val: string(setter.Val[:len(setter.Val)-1])},
					f:   tks[0],
					fv:  v.Val,
					s:   setter,
					sv:  val,
				})
		}
	}
	p.vars[j] = v
}
