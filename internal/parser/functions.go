package parser

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/fract-lang/fract/internal/functions/built_in"
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
	// Is built-in function?
	if c.f.Tks == nil {
		// Add name token for exceptions.
		c.f.Tks = []obj.Tokens{{c.name}}
		switch c.f.Name {
		case "print":
			built_in.Print(c.f, c.args)
		case "input":
			return built_in.Input(c.f, c.args)
		case "len":
			return built_in.Len(c.f, c.args)
		case "range":
			return built_in.Range(c.f, c.args)
		case "calloc":
			return built_in.Calloc(c.f, c.args)
		case "realloc":
			return built_in.Realloc(c.f, c.args)
		case "string":
			return built_in.String(c.f, c.args)
		case "int":
			return built_in.Int(c.f, c.args)
		case "float":
			return built_in.Float(c.f, c.args)
		case "append":
			return built_in.Append(c.f, c.args)
		default:
			built_in.Exit(c.f, c.args)
		}
		return retv
	}
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
				} else if c.src.process(tks) == fract.FUNCReturn {
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
	if b.P.M != "" {
		defers = defers[:dlen]
		panic(b.P.M)
	}
	for i := len(defers) - 1; i >= dlen; i-- {
		defers[i].call()
	}
	defers = defers[:dlen]
	return retv
}

// isParamSet Argument type is param set?
func isParamSet(tks obj.Tokens) bool {
	return tks[0].T == fract.Name && tks[1].Val == "="
}

// paramsArgVals decompose and returns params values.
func (p *Parser) paramsArgVals(tks obj.Tokens, i, lstComma *int) obj.Value {
	retv := obj.Value{D: []obj.Data{}, Arr: true}
	bc := 0
	for ; *i < len(tks); *i++ {
		tk := tks[*i]
		if tk.T == fract.Brace {
			if tk.Val == "(" || tk.Val == "{" || tk.Val == "[" {
				bc++
			} else {
				bc--
			}
		} else if tk.T == fract.Comma && bc == 0 {
			vtks := tks.Sub(*lstComma, *i-*lstComma)
			if isParamSet(*vtks) {
				*i -= 4
				return retv
			}
			v := p.procVal(*vtks)
			if v.Arr {
				retv.D = append(retv.D, obj.Data{D: v.D, T: obj.VArray})
			} else {
				retv.D = append(retv.D, v.D...)
			}
			*lstComma = *i + 1
		}
	}
	if *lstComma < len(tks) {
		vtks := tks[*lstComma:]
		if isParamSet(vtks) {
			*i -= 4
			return retv
		}
		v := p.procVal(vtks)
		if v.Arr {
			retv.D = append(retv.D, obj.Data{D: v.D, T: obj.VArray})
		} else {
			retv.D = append(retv.D, v.D...)
		}
	}
	return retv
}

type funcArgInfo struct {
	f        obj.Func
	names    *[]string
	tks      obj.Tokens
	tk       obj.Token
	index    *int
	count    *int
	lstComma *int
}

// Process function argument.
func (p *Parser) procFuncArg(i funcArgInfo) obj.Var {
	var paramSet bool
	l := *i.index - *i.lstComma
	if l < 1 {
		fract.IPanic(i.tk, obj.SyntaxPanic, "Value is not given!")
	} else if *i.count >= len(i.f.Params) {
		fract.IPanic(i.tk, obj.SyntaxPanic, "Argument overflow!")
	}
	param := i.f.Params[*i.count]
	v := obj.Var{Name: param.Name}
	vtks := *i.tks.Sub(*i.lstComma, l)
	i.tk = vtks[0]
	// Check param set.
	if l >= 2 && isParamSet(vtks) {
		l -= 2
		if l < 1 {
			fract.IPanic(i.tk, obj.SyntaxPanic, "Value is not given!")
		}
		for _, pr := range i.f.Params {
			if pr.Name == i.tk.Val {
				for _, name := range *i.names {
					if name == i.tk.Val {
						fract.IPanic(i.tk, obj.SyntaxPanic, "Keyword argument repeated!")
					}
				}
				*i.count++
				paramSet = true
				*i.names = append(*i.names, i.tk.Val)
				retv := obj.Var{Name: i.tk.Val}
				//Parameter is params typed?
				if pr.Params {
					*i.lstComma += 2
					retv.Val = p.paramsArgVals(i.tks, i.index, i.lstComma)
				} else {
					retv.Val = p.procVal(vtks[2:])
				}
				return retv
			}
		}
		fract.IPanic(i.tk, obj.NamePanic, "Parameter is not defined in this name: "+i.tk.Val)
	}
	if paramSet {
		fract.IPanic(i.tk, obj.SyntaxPanic, "After the parameter has been given a special value, all parameters must be shown privately!")
	}
	*i.count++
	*i.names = append(*i.names, v.Name)
	// Parameter is params typed?
	if param.Params {
		v.Val = p.paramsArgVals(i.tks, i.index, i.lstComma)
	} else {
		v.Val = p.procVal(vtks)
	}
	return v
}

// Process function call model and initialize model instance.
func (p *Parser) funcCallModel(tks obj.Tokens) funcCall {
	name := tks[0]
	// Name is not defined?
	namei, src := p.funcIndexByName(name)
	var f obj.Func
	if namei == -1 {
		name := name
		if j := strings.IndexByte(name.Val, '.'); j != -1 {
			if p.importIndexByName(name.Val[:j]) == -1 {
				fract.IPanic(name, obj.NamePanic, "'"+name.Val[:j]+"' is not defined!")
			}
			src = p.Imports[p.importIndexByName(name.Val[:j])].Src
			name.Val = name.Val[j+1:]
			for _, v := range src.vars {
				if unicode.IsUpper(rune(v.Name[0])) && v.Name == name.Val && !v.Val.Arr && v.Val.D[0].T == obj.VFunc {
					name.F = nil
					f = v.Val.D[0].D.(obj.Func)
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
			fract.IPanic(name, obj.NamePanic, "Function is not defined in this name: "+name.Val)
		}
	} else {
		f = src.funcs[namei]
	}
	var (
		names []string
		args  []obj.Var
		count = 0
	)
	// Decompose arguments.
	if tks, _ = decomposeBrace(&tks, "(", ")", false); tks != nil {
		var (
			inf = funcArgInfo{
				f:        f,
				names:    &names,
				tks:      tks,
				count:    &count,
				index:    new(int),
				lstComma: new(int),
			}
			bc = 0
		)
		for *inf.index = 0; *inf.index < len(tks); *inf.index++ {
			inf.tk = tks[*inf.index]
			if inf.tk.T == fract.Brace {
				if inf.tk.Val == "(" || inf.tk.Val == "{" || inf.tk.Val == "[" {
					bc++
				} else {
					bc--
				}
			} else if inf.tk.T == fract.Comma && bc == 0 {
				args = append(args, p.procFuncArg(inf))
				*inf.lstComma = *inf.index + 1
			}
		}
		if *inf.lstComma < len(tks) {
			inf.tk = tks[*inf.lstComma]
			tkslen := len(tks)
			inf.index = &tkslen
			args = append(args, p.procFuncArg(inf))
		}
	}
	// All parameters is not defined?
	if count < len(f.Params)-f.DefaultParamCount {
		var sb strings.Builder
		sb.WriteString("All required positional arguments is not given:")
		for _, p := range f.Params {
			if p.Default.D != nil {
				break
			}
			msg := " '" + p.Name + "',"
			for _, name := range names {
				if p.Name == name {
					msg = ""
					break
				}
			}
			sb.WriteString(msg)
		}
		fract.IPanic(name, obj.PlainPanic, sb.String()[:sb.Len()-1])
	}
	// Check default values.
	for ; count < len(f.Params); count++ {
		p := f.Params[count]
		if p.Default.D != nil {
			args = append(args, obj.Var{Name: p.Name, Val: p.Default})
		}
	}
	return funcCall{
		f:    f,
		name: name,
		src:  src,
		args: args,
	}
}

// funcCall call function and returns returned value.
func (p *Parser) funcCall(tks obj.Tokens) obj.Value {
	return p.funcCallModel(tks).call()
}

// Process function declaration.
func (p *Parser) funcdec(tks obj.Tokens, protected bool) {
	tkslen := len(tks)
	name := tks[1]
	// Name is not name?
	if name.T != fract.Name {
		fract.IPanic(name, obj.SyntaxPanic, "Invalid name!")
	} else if strings.Contains(name.Val, ".") {
		fract.IPanic(name, obj.SyntaxPanic, "Names is cannot include dot!")
	}
	// Name is already defined?
	if line := p.definedName(name); line != -1 {
		fract.IPanic(name, obj.NamePanic, "\""+name.Val+"\" is already defined at line: "+fmt.Sprint(line))
	}
	// Function parentheses are not defined?
	if tkslen < 4 {
		fract.IPanic(name, obj.SyntaxPanic, "Function parentheses is not found!")
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
		fract.IPanic(dtToken, obj.SyntaxPanic, "Invalid syntax!")
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
					fract.IPanic(pr, obj.SyntaxPanic, "Parameter name is not found!")
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
						fract.IPanic(ptks[start-1], obj.SyntaxPanic, "Value is not given!")
					}
					lstp.Default = p.procVal(*ptks.Sub(start, i-start))
					if lstp.Params && !lstp.Default.Arr {
						fract.IPanic(pr, obj.ValuePanic, "Params parameter is can only take array values!")
					}
					f.Params[len(f.Params)-1] = lstp
					f.DefaultParamCount++
					defaultDef = true
					continue
				}
				if lstp.Default.D == nil && defaultDef {
					fract.IPanic(pr, obj.SyntaxPanic, "All parameters after a given parameter with a default value must take a default value!")
				} else if pr.T != fract.Comma {
					fract.IPanic(pr, obj.SyntaxPanic, "Comma is not found!")
				}
			}
		}
		if lstp.Default.D == nil && defaultDef {
			fract.IPanic(tks[len(tks)-1], obj.SyntaxPanic, "All parameters after a given parameter with a default value must take a default value!")
		}
	}
	p.skipBlock(false)
	f.Tks = p.Tks[f.Ln : p.i+1]
	f.Ln = name.Ln
	p.funcs = append(p.funcs, f)
}
