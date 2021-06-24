package parser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/fract-lang/fract/internal/functions/embed"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// Instance for function calls.
type funcCall struct {
	f    obj.Func
	name obj.Token
	src  *Parser
	args []obj.Var
}

func (c funcCall) call() obj.Value {
	var retv obj.Value
	// Is embed function?
	if c.f.Tks == nil {
		// Add name token for exceptions.
		c.f.Tks = []obj.Tokens{{c.name}}
		switch c.f.Name {
		case "print":
			embed.Print(c.f, c.args)
		case "input":
			return embed.Input(c.f, c.args)
		case "len":
			return embed.Len(c.f, c.args)
		case "range":
			return embed.Range(c.f, c.args)
		case "make":
			return embed.Make(c.f, c.args)
		case "string":
			return embed.String(c.f, c.args)
		case "int":
			return embed.Int(c.f, c.args)
		case "float":
			return embed.Float(c.f, c.args)
		case "append":
			return embed.Append(c.f, c.args)
		default:
			embed.Exit(c.f, c.args)
		}
	} else {
		// Process block.
		vars := c.src.vars
		dlen := len(defers)
		old := c.src.funcTempVars
		if c.src.funcTempVars == -1 {
			c.src.funcTempVars = 0
		}
		if c.src.funcTempVars == 0 {
			c.src.vars = append(c.args, c.src.vars...)
		} else {
			c.src.vars = append(c.args, c.src.vars[:len(c.src.vars)-c.src.funcTempVars]...)
		}
		c.src.funcCount++
		c.src.funcTempVars = len(c.args)
		flen := len(c.src.funcs)
		namei := c.src.i
		itks := c.src.Tks
		c.src.Tks = c.f.Tks
		c.src.i = -1
		// Interpret block.
		b := obj.Block{
			Try: func() {
				for {
					c.src.i++
					tks := c.src.Tks[c.src.i]
					if tks[0].T == fract.End { // Block is ended.
						break
					} else if c.src.processTokens(tks) == fract.FUNCReturn {
						if c.src.retVal == nil {
							break
						}
						retv = *c.src.retVal
						c.src.retVal = nil
						break
					}
				}
			},
		}
		b.Do()
		c.src.Tks = itks
		// Remove temporary functions.
		c.src.funcs = c.src.funcs[:flen]
		// Remove temporary variables.
		c.src.vars = vars
		c.src.funcCount--
		c.src.funcTempVars = old
		c.src.i = namei
		if b.E != nil {
			defers = defers[:dlen]
			panic(fmt.Errorf(b.E.Msg))
		}
		for i := len(defers) - 1; i >= dlen; i-- {
			defers[i].call()
		}
		defers = defers[:dlen]
	}
	return retv
}

// isParamSet Argument type is param set?
func isParamSet(tks obj.Tokens) bool { return tks[0].T == fract.Name && tks[1].Val == "=" }

// getParamsArgumentValues decompose and returns params values.
func (p *Parser) getParamsArgumentValues(tks obj.Tokens, pos, bc, lstComma *int) obj.Value {
	retv := obj.Value{D: []obj.Data{}, Arr: true}
	for ; *pos < len(tks); *pos++ {
		tk := tks[*pos]
		if tk.T == fract.Brace {
			if tk.Val == "(" || tk.Val == "{" || tk.Val == "[" {
				*bc++
			} else {
				*bc--
			}
		} else if tk.T == fract.Comma && *bc == 0 {
			vtks := tks.Sub(*lstComma, *pos-*lstComma)
			if isParamSet(*vtks) {
				*pos -= 4
				return retv
			}
			v := p.processValue(*vtks)
			if v.Arr {
				retv.D = append(retv.D, obj.Data{D: v.D, T: obj.VArray})
			} else {
				retv.D = append(retv.D, v.D...)
			}
			*lstComma = *pos + 1
		}
	}
	if *lstComma < len(tks) {
		vtks := tks[*lstComma:]
		if isParamSet(vtks) {
			*pos -= 4
			return retv
		}
		v := p.processValue(vtks)
		if v.Arr {
			retv.D = append(retv.D, obj.Data{D: v.D, T: obj.VArray})
		} else {
			retv.D = append(retv.D, v.D...)
		}
	}
	return retv
}

func (p *Parser) processArgument(f obj.Func, names *[]string, tks obj.Tokens, tk obj.Token, index, count, bc, lstComma *int) obj.Var {
	var paramSet bool
	l := *index - *lstComma
	if l < 1 {
		fract.Error(tk, "Value is not defined!")
	} else if *count >= len(f.Params) {
		fract.Error(tk, "Argument overflow!")
	}
	param := f.Params[*count]
	v := obj.Var{Name: param.Name}
	vtks := *tks.Sub(*lstComma, l)
	tk = vtks[0]
	// Check param set.
	if l >= 2 && isParamSet(vtks) {
		l -= 2
		if l < 1 {
			fract.Error(tk, "Value is not defined!")
		}
		for _, pr := range f.Params {
			if pr.Name == tk.Val {
				for _, name := range *names {
					if name == tk.Val {
						fract.Error(tk, "Keyword argument repeated!")
					}
				}
				*count++
				paramSet = true
				*names = append(*names, tk.Val)
				retv := obj.Var{Name: tk.Val}
				//Parameter is params typed?
				if pr.Params {
					*lstComma += 2
					retv.Val = p.getParamsArgumentValues(tks, index, bc, lstComma)
				} else {
					retv.Val = p.processValue(vtks[2:])
				}
				return retv
			}
		}
		fract.Error(tk, "Parameter is not defined in this name: "+tk.Val)
	}
	if paramSet {
		fract.Error(tk, "After the parameter has been given a special value, all parameters must be shown privately!")
	}
	*count++
	*names = append(*names, v.Name)
	// Parameter is params typed?
	if param.Params {
		v.Val = p.getParamsArgumentValues(tks, index, bc, lstComma)
	} else {
		v.Val = p.processValue(vtks)
	}
	return v
}

// Process function call model and initialize moden instance.
func (p *Parser) processFunctionCallModel(tks obj.Tokens) funcCall {
	name := tks[0]
	// Name is not defined?
	namei, src := p.functionIndexByName(name)
	var f obj.Func
	if namei == -1 {
		name := name
		if j := strings.Index(name.Val, "."); j != -1 {
			if p.importIndexByName(name.Val[:j]) == -1 {
				fract.Error(name, "'"+name.Val[:j]+"' is not defined!")
			}
			src = p.Imports[p.importIndexByName(name.Val[:j])].Src
			name.Val = name.Val[j+1:]
			for _, current := range src.vars {
				if unicode.IsUpper(rune(current.Name[0])) && current.Name == name.Val && !current.Val.Arr && current.Val.D[0].T == obj.VFunc {
					name.F = nil
					f = current.Val.D[0].D.(obj.Func)
					break
				}
			}
		} else {
			for _, current := range p.vars {
				if current.Name == name.Val && !current.Val.Arr && current.Val.D[0].T == obj.VFunc {
					name.F = nil
					f = current.Val.D[0].D.(obj.Func)
					src = p
					break
				}
			}
		}
		if name.F != nil {
			fract.Error(name, "Function is not defined in this name: "+name.Val)
		}
	} else {
		f = src.funcs[namei]
	}
	var (
		names = new([]string)
		count = new(int)
		args  []obj.Var
	)
	// Decompose arguments.
	if tks, _ = decomposeBrace(&tks, "(", ")", false); tks != nil {
		var (
			bc       = new(int)
			lstComma = new(int)
		)
		for i := 0; i < len(tks); i++ {
			tk := tks[i]
			if tk.T == fract.Brace {
				if tk.Val == "(" || tk.Val == "{" || tk.Val == "[" {
					*bc++
				} else {
					*bc--
				}
			} else if tk.T == fract.Comma && *bc == 0 {
				args = append(args, p.processArgument(f, names, tks, tk, &i, count, bc, lstComma))
				*lstComma = i + 1
			}
		}
		if *lstComma < len(tks) {
			tkslen := len(tks)
			args = append(args, p.processArgument(f, names, tks, tks[*lstComma], &tkslen, count, bc, lstComma))
		}
	}
	// All parameters is not defined?
	if *count < len(f.Params)-f.DefaultParamCount {
		var sb strings.Builder
		sb.WriteString("All required positional parameters is not defined:")
		for _, p := range f.Params {
			if p.Default.D != nil {
				break
			}
			msg := " '" + p.Name + "',"
			for _, name := range *names {
				if p.Name == name {
					msg = ""
					break
				}
			}
			sb.WriteString(msg)
		}
		fract.Error(name, sb.String()[:sb.Len()-1])
	}
	// Check default values.
	for ; *count < len(f.Params); *count++ {
		p := f.Params[*count]
		if p.Default.D != nil {
			args = append(args,
				obj.Var{
					Name: p.Name,
					Val:  p.Default,
				})
		}
	}
	return funcCall{
		f:    f,
		name: name,
		src:  src,
		args: args,
	}
}

// processFunctionCall call function and returns returned value.
func (p *Parser) processFunctionCall(tks obj.Tokens) obj.Value {
	return p.processFunctionCallModel(tks).call()
}

func (p *Parser) processFunctionDeclaration(tks obj.Tokens, protected bool) {
	tkslen := len(tks)
	name := tks[1]
	// Name is not name?
	if name.T != fract.Name {
		fract.Error(name, "This is not a valid name!")
	} else if strings.Contains(name.Val, ".") {
		fract.Error(name, "Names is cannot include dot!")
	}
	// Name is already defined?
	if line := p.definedName(name); line != -1 {
		fract.Error(name, "\""+name.Val+"\" is already defined at line: "+fmt.Sprint(line))
	}
	// Function parentheses are not defined?
	if tkslen < 4 {
		fract.Error(name, "Where is the function parentheses?")
	}
	p.i++
	f := obj.Func{
		Name:      name.Val,
		Ln:        p.i,
		Params:    []obj.Param{},
		Protected: protected,
	}
	dtToken := tks[tkslen-1]
	if dtToken.T != fract.Brace || dtToken.Val != ")" {
		fract.Error(dtToken, "Invalid syntax!")
	}
	if paramtks := tks.Sub(3, tkslen-4); paramtks != nil {
		ptks := *paramtks
		// Decompose function parameters.
		pname, defaultDef := true, false
		var lstp obj.Param
		for i := 0; i < len(ptks); i++ {
			pr := ptks[i]
			if pname {
				if pr.T == fract.Params {
					continue
				} else if pr.T != fract.Name {
					fract.Error(pr, "Parameter name is not found!")
				}
				lstp = obj.Param{
					Name:   pr.Val,
					Params: i > 0 && ptks[i-1].T == fract.Params,
				}
				f.Params = append(f.Params, lstp)
				pname = false
				continue
			} else {
				pname = true
				// Default value definition?
				if pr.Val == "=" {
					bc := 0
					i++
					start := i
					for ; i < len(ptks); i++ {
						pr = ptks[i]
						if pr.T == fract.Brace {
							if pr.Val == "{" || pr.Val == "(" || pr.Val == "[" {
								bc++
							} else {
								bc--
							}
						} else if pr.T == fract.Comma {
							break
						}
					}
					if i-start < 1 {
						fract.Error(ptks[start-1], "Value is not defined!")
					}
					lstp.Default = p.processValue(*ptks.Sub(start, i-start))
					if lstp.Params && !lstp.Default.Arr {
						fract.Error(pr, "Params parameter is can only take array values!")
					}
					f.Params[len(f.Params)-1] = lstp
					f.DefaultParamCount++
					defaultDef = true
					continue
				}
				if lstp.Default.D == nil && defaultDef {
					fract.Error(pr, "All parameters after a given parameter with a default value must take a default value!")
				} else if pr.T != fract.Comma {
					fract.Error(pr, "Comma is not found!")
				}
			}
		}
		if lstp.Default.D == nil && defaultDef {
			fract.Error(tks[len(tks)-1], "All parameters after a given parameter with a default value must take a default value!")
		}
	}
	p.skipBlock(false)
	f.Tks = p.Tks[f.Ln : p.i+1]
	f.Ln = name.Ln
	p.funcs = append(p.funcs, f)
}
