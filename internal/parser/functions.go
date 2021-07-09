package parser

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/internal/functions/built_in"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

// Func instance.
type Func struct {
	name              string
	src               *Parser
	ln                int          // Line of define.
	tks               []obj.Tokens // Block content of function.
	params            []obj.Param
	defaultParamCount int
	protected         bool
}

// Instance for function calls.
type funcCall struct {
	f     Func
	errTk obj.Token
	args  []obj.Var
}

func (c funcCall) call() obj.Value {
	var retv obj.Value
	// Is built-in function?
	if c.f.tks == nil {
		switch c.f.name {
		case "print":
			built_in.Print(c.errTk, c.args)
		case "input":
			return built_in.Input(c.args)
		case "len":
			return built_in.Len(c.args)
		case "range":
			return built_in.Range(c.errTk, c.args)
		case "calloc":
			return built_in.Calloc(c.errTk, c.args)
		case "realloc":
			return built_in.Realloc(c.errTk, c.args)
		case "string":
			return built_in.String(c.args)
		case "int":
			return built_in.Int(c.args)
		case "float":
			return built_in.Float(c.args)
		case "append":
			return built_in.Append(c.errTk, c.args)
		default:
			built_in.Exit(c.errTk, c.args)
		}
		return retv
	}
	// Process block.
	vars := c.f.src.vars
	dlen := len(defers)
	old := c.f.src.funcTempVars
	if c.f.src.funcTempVars == -1 {
		c.f.src.funcTempVars = 0
	}
	if c.f.src.funcTempVars == 0 {
		c.f.src.vars = append(c.args, c.f.src.vars...)
	} else {
		c.f.src.vars = append(c.args, c.f.src.vars[:len(c.f.src.vars)-c.f.src.funcTempVars]...)
	}
	c.f.src.funcCount++
	c.f.src.funcTempVars = len(c.args)
	flen := len(c.f.src.funcs)
	namei := c.f.src.i
	itks := c.f.src.Tks
	c.f.src.Tks = c.f.tks
	c.f.src.i = -1
	// Interpret block.
	b := obj.Block{
		Try: func() {
			for {
				c.f.src.i++
				tks := c.f.src.Tks[c.f.src.i]
				if tks[0].T == fract.End { // Block is ended.
					break
				} else if c.f.src.process(tks) == fract.FUNCReturn {
					if c.f.src.retVal == nil {
						break
					}
					retv = *c.f.src.retVal
					c.f.src.retVal = nil
					break
				}
			}
		},
	}
	b.Do()
	c.f.src.Tks = itks
	// Remove temporary functions.
	c.f.src.funcs = c.f.src.funcs[:flen]
	// Remove temporary variables.
	c.f.src.vars = vars
	c.f.src.funcCount--
	c.f.src.funcTempVars = old
	c.f.src.i = namei
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
	return tks[0].T == fract.Name && tks[1].V == "="
}

// paramsArgVals decompose and returns params values.
func (p *Parser) paramsArgVals(tks obj.Tokens, i, lstComma *int) obj.Value {
	retv := obj.Value{D: []obj.Data{}, Arr: true}
	bc := 0
	for ; *i < len(tks); *i++ {
		switch tk := tks[*i]; tk.T {
		case fract.Brace:
			switch tk.V {
			case "{", "[", "(":
				bc++
			default:
				bc--
			}
		case fract.Comma:
			if bc != 0 {
				break
			}
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
	f        Func
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
	} else if *i.count >= len(i.f.params) {
		fract.IPanic(i.tk, obj.SyntaxPanic, "Argument overflow!")
	}
	param := i.f.params[*i.count]
	v := obj.Var{Name: param.Name}
	vtks := *i.tks.Sub(*i.lstComma, l)
	i.tk = vtks[0]
	// Check param set.
	if l >= 2 && isParamSet(vtks) {
		l -= 2
		if l < 1 {
			fract.IPanic(i.tk, obj.SyntaxPanic, "Value is not given!")
		}
		for _, pr := range i.f.params {
			if pr.Name == i.tk.V {
				for _, name := range *i.names {
					if name == i.tk.V {
						fract.IPanic(i.tk, obj.SyntaxPanic, "Keyword argument repeated!")
					}
				}
				*i.count++
				paramSet = true
				*i.names = append(*i.names, i.tk.V)
				retv := obj.Var{Name: i.tk.V}
				//Parameter is params typed?
				if pr.Params {
					*i.lstComma += 2
					retv.V = p.paramsArgVals(i.tks, i.index, i.lstComma)
				} else {
					retv.V = p.procVal(vtks[2:])
				}
				return retv
			}
		}
		fract.IPanic(i.tk, obj.NamePanic, "Parameter is not defined in this name: "+i.tk.V)
	}
	if paramSet {
		fract.IPanic(i.tk, obj.SyntaxPanic, "After the parameter has been given a special value, all parameters must be shown privately!")
	}
	*i.count++
	*i.names = append(*i.names, v.Name)
	// Parameter is params typed?
	if param.Params {
		v.V = p.paramsArgVals(i.tks, i.index, i.lstComma)
	} else {
		v.V = p.procVal(vtks)
	}
	return v
}

// Process function call model and initialize model instance.
func (p *Parser) funcCallModel(f Func, tks obj.Tokens) funcCall {
	var (
		names []string
		args  []obj.Var
		count = 0
		tk    = tks[0]
	)
	// Decompose arguments.
	if tks, _ = decomposeBrace(&tks, "(", ")"); tks != nil {
		var (
			inf = funcArgInfo{
				f:        f,
				names:    &names,
				tk:       tk,
				tks:      tks,
				count:    &count,
				index:    new(int),
				lstComma: new(int),
			}
			bc = 0
		)
		for *inf.index = 0; *inf.index < len(tks); *inf.index++ {
			switch inf.tk = tks[*inf.index]; inf.tk.T {
			case fract.Brace:
				switch inf.tk.V {
				case "{", "[", "(":
					bc++
				default:
					bc--
				}
			case fract.Comma:
				if bc != 0 {
					break
				}
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
	if count < len(f.params)-f.defaultParamCount {
		var sb strings.Builder
		sb.WriteString("All required positional arguments is not given:")
		for _, p := range f.params {
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
		fract.IPanic(tk, obj.PlainPanic, sb.String()[:sb.Len()-1])
	}
	// Check default values.
	for ; count < len(f.params); count++ {
		p := f.params[count]
		if p.Default.D != nil {
			args = append(args, obj.Var{Name: p.Name, V: p.Default})
		}
	}
	return funcCall{
		f:     f,
		errTk: tk,
		args:  args,
	}
}

// funcCall call function and returns returned value.
func (p *Parser) funcCall(f Func, tks obj.Tokens) obj.Value {
	return p.funcCallModel(f, tks).call()
}

// Process function declaration.
func (p *Parser) funcdec(tks obj.Tokens, protected bool) {
	tkslen := len(tks)
	name := tks[1]
	// Name is not name?
	if name.T != fract.Name {
		fract.IPanic(name, obj.SyntaxPanic, "Invalid name!")
	} else if strings.Contains(name.V, ".") {
		fract.IPanic(name, obj.SyntaxPanic, "Names is cannot include dot!")
	}
	// Name is already defined?
	if line := p.definedName(name); line != -1 {
		fract.IPanic(name, obj.NamePanic, "\""+name.V+"\" is already defined at line: "+fmt.Sprint(line))
	}
	// Function parentheses are not defined?
	if tkslen < 4 {
		fract.IPanic(name, obj.SyntaxPanic, "Function parentheses is not found!")
	}
	p.i++
	f := Func{
		name:      name.V,
		ln:        p.i,
		params:    []obj.Param{},
		protected: protected,
		src:       p,
	}
	dtToken := tks[tkslen-1]
	if dtToken.T != fract.Brace || dtToken.V != ")" {
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
				switch pr.T {
				case fract.Params:
					continue
				case fract.Name:
					fract.IPanic(pr, obj.SyntaxPanic, "Parameter name is not found!")
				}
				lstp = obj.Param{Name: pr.V, Params: i > 0 && ptks[i-1].T == fract.Params}
				f.params = append(f.params, lstp)
				pname = false
				continue
			} else {
				pname = true
				// Default value definition?
				if pr.V == "=" {
					bc := 0
					i++
					start := i
					for ; i < len(ptks); i++ {
						pr = ptks[i]
						if pr.T == fract.Brace {
							switch pr.V {
							case "{", "[", "(":
								bc++
							default:
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
					f.params[len(f.params)-1] = lstp
					f.defaultParamCount++
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
	f.tks = p.Tks[f.ln : p.i+1]
	f.ln = name.Ln
	p.funcs = append(p.funcs, f)
}
