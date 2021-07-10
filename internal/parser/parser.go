package parser

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"unicode"

	"github.com/fract-lang/fract/internal/lex"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/obj"
)

var (
	defers []funcCall
)

// Parser of Fract.
type Parser struct {
	vars         []obj.Var
	funcs        []function
	funcTempVars int // Count of function temporary variables.
	loopCount    int
	funcCount    int
	i            int
	retVal       *obj.Value // Pointer of last return value.

	L       *lex.Lex
	Tks     []obj.Tokens // All Tokens of code file.
	Imports []importInfo
}

// New returns instance of parser related to file.
func New(fp string) *Parser {
	f, _ := os.Open(fp)
	bytes, _ := os.ReadFile(fp)
	sf := &obj.File{P: fp, F: f}
	if runtime.GOOS == "linux" {
		sf.Lns = strings.Split(string(bytes), "\r\n")
	} else {
		sf.Lns = strings.Split(string(bytes), "\n")
	}
	return &Parser{
		funcTempVars: -1,
		L:            &lex.Lex{F: sf, Ln: 1},
	}
}

// NewStdin returns new instance of parser from standard input.
func NewStdin() *Parser {
	return &Parser{
		funcTempVars: -1,
		L: &lex.Lex{
			F:  &obj.File{P: "<stdin>"},
			Ln: 1,
		},
	}
}

// ready interpreter to process.
func (p *Parser) ready() {
	/// Tokenize all lines.
	for !p.L.Fin {
		if ctks := p.L.Next(); ctks != nil {
			p.Tks = append(p.Tks, ctks)
		}
	}
	// Change blocks.
	bc := 0
	mbc := 0
	lst := -1
	for i, tks := range p.Tks {
		if fst := tks[0]; fst.T == fract.End {
			if len(tks) > 1 {
				fract.IPanic(tks[1], obj.SyntaxPanic, "Invalid syntax!")
			}
			bc--
			if bc < 0 {
				fract.IPanic(fst, obj.SyntaxPanic, "The extra block end defined!")
			}
		} else if fst.T == fract.Macro {
			if IsBlock(tks) {
				mbc++
				if mbc == 1 {
					lst = i
				}
			} else if tks[1].T == fract.End {
				if len(tks) > 2 {
					fract.IPanic(tks[2], obj.SyntaxPanic, "Invalid syntax!")
				}
				mbc--
				if mbc < 0 {
					fract.IPanic(fst, obj.SyntaxPanic, "The extra block end defined!")
				}
			}
		} else if IsBlock(tks) {
			bc++
			if bc == 1 {
				lst = i
			}
		}
	}
	if bc > 0 || mbc > 0 { // Check blocks.
		fract.IPanic(p.Tks[lst][0], obj.SyntaxPanic, "Block is expected ending...")
	}
}

func (p *Parser) importLocal() {
	dir, _ := os.Getwd()
	if pdir := path.Dir(p.L.F.P); pdir != "." {
		dir = path.Join(dir, pdir)
	}
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	_, mainName := filepath.Split(p.L.F.P)
	for _, info := range infos {
		// Skip directories.
		if info.IsDir() || !strings.HasSuffix(info.Name(), fract.Ext) || info.Name() == mainName {
			continue
		}
		src := New(path.Join(dir, info.Name()))
		src.loopCount = -1 //! Tag as import source.
		src.Import()
		p.funcs = append(p.funcs, src.funcs...)
		p.vars = append(p.vars, src.vars...)
		p.Imports = append(p.Imports, src.Imports...)
	}
}

func (p *Parser) Interpret() {
	if p.L.F.P == "<stdin>" {
		// Interpret all lines.
		for p.i = 0; p.i < len(p.Tks); p.i++ {
			p.process(p.Tks[p.i])
		}
		goto end
	}
	// Lexer is finished.
	if p.L.Fin {
		return
	}
	p.ready()
	p.importLocal()
	// Interpret all lines.
	for p.i = 0; p.i < len(p.Tks); p.i++ {
		p.process(p.Tks[p.i])
	}
end:
	for i := len(defers) - 1; i >= 0; i-- {
		defers[i].call()
	}
}

// Process array indexes for access to elements.
func indexes(arr, val obj.Value, tk obj.Token) []int {
	if val.Arr {
		var i []int
		for _, d := range val.D {
			if d.T != obj.VInt {
				fract.IPanic(tk, obj.ValuePanic, "Only integer values can used in index access!")
			}
			pos, err := strconv.Atoi(d.String())
			if err != nil {
				fract.IPanic(tk, obj.OutOfRangePanic, "Value out of range!")
			}
			if arr.Arr {
				pos = procIndex(len(arr.D), pos)
			} else {
				pos = procIndex(len(arr.D[0].String()), pos)
			}
			if pos == -1 {
				fract.IPanic(tk, obj.OutOfRangePanic, "Index is out of range!")
			}
			i = append(i, pos)
		}
		return i
	}
	if val.D[0].T != obj.VInt {
		fract.IPanic(tk, obj.ValuePanic, "Only integer values can used in index access!")
	}
	pos, err := strconv.Atoi(val.String())
	if err != nil {
		fract.IPanic(tk, obj.OutOfRangePanic, "Value out of range!")
	}
	if arr.Arr {
		pos = procIndex(len(arr.D), pos)
	} else {
		pos = procIndex(len(arr.D[0].String()), pos)
	}
	if pos == -1 {
		fract.IPanic(tk, obj.OutOfRangePanic, "Index is out of range!")
	}
	return []int{pos}
}

// skipBlock skip to block end.
func (p *Parser) skipBlock(blk bool) {
	if blk {
		if IsBlock(p.Tks[p.i]) {
			p.i++
		} else {
			return
		}
	}
	c := 1
	p.i--
	for {
		p.i++
		tks := p.Tks[p.i]
		if fst := tks[0]; fst.T == fract.End {
			c--
		} else if fst.T == fract.Macro {
			if IsBlock(tks) {
				c++
			} else if tks[1].T == fract.End {
				c--
			}
		} else if IsBlock(tks) {
			c++
		}
		if c == 0 {
			return
		}
	}
}

// TYPES
// 'f' -> Function.
// 'v' -> Variable.
// Returns define by name.
func (p *Parser) defByName(name obj.Token) (int, rune, *Parser) {
	pos, src := p.funcIndexByName(name)
	if pos != -1 {
		return pos, 'f', src
	}
	pos, src = p.varIndexByName(name)
	if pos != -1 {
		return pos, 'v', src
	}
	return -1, '-', nil
}

// Returns index of name is exist name, returns -1 if not.
func (p *Parser) definedName(name obj.Token) int {
	if name.V[0] == '-' { // Ignore minus.
		name.V = name.V[1:]
	}
	for _, f := range p.funcs {
		if f.name == name.V {
			return f.ln
		}
	}
	for _, v := range p.vars {
		if v.Name == name.V {
			return v.Ln
		}
	}
	return -1
}

//! This code block very like to varIndexByName function.
//! If you change here, probably you must change there too.

// funcIndexByName returns index of function by name.
func (p *Parser) funcIndexByName(name obj.Token) (int, *Parser) {
	if name.V[0] == '-' { // Ignore minus.
		name.V = name.V[1:]
	}
	if i := strings.IndexByte(name.V, '.'); i != -1 {
		if p.importIndexByName(name.V[:i]) == -1 {
			fract.IPanic(name, obj.NamePanic, "'"+name.V[:i]+"' is not defined!")
		}
		p = p.Imports[p.importIndexByName(name.V[:i])].Src
		name.V = name.V[i+1:]
		for i, current := range p.funcs {
			if (current.tks == nil || unicode.IsUpper(rune(current.name[0]))) && current.name == name.V {
				return i, p
			}
		}
		return -1, nil
	}
	for j, f := range p.funcs {
		if f.name == name.V {
			return j, p
		}
	}
	return -1, nil
}

//! This code block very like to funcIndexByName function.
//! If you change here, probably you must change there too.

// varIndexByName returns index of variable by name.
func (p *Parser) varIndexByName(name obj.Token) (int, *Parser) {
	if name.V[0] == '-' { // Ignore minus.
		name.V = name.V[1:]
	}
	if i := strings.IndexByte(name.V, '.'); i != -1 {
		if iindex := p.importIndexByName(name.V[:i]); iindex == -1 {
			fract.IPanic(name, obj.NamePanic, "'"+name.V[:i]+"' is not defined!")
		} else {
			p = p.Imports[iindex].Src
		}
		name.V = name.V[i+1:]
		for i, v := range p.vars {
			if (v.Ln == -1 || unicode.IsUpper(rune(v.Name[0]))) && v.Name == name.V {
				return i, p
			}
		}
		return -1, nil
	}
	for j, v := range p.vars {
		if v.Name == name.V {
			return j, p
		}
	}
	return -1, nil
}

// importIndexByName returns index of import by name.
func (p *Parser) importIndexByName(name string) int {
	for i, imp := range p.Imports {
		if imp.Name == name {
			return i
		}
	}
	return -1
}

// Check arithmetic processes validity.
func arithmeticProcesses(tks obj.Tokens) []obj.Tokens {
	var (
		procs []obj.Tokens
		part  obj.Tokens
		opr   bool
		b     int
	)
	for i := 0; i < len(tks); i++ {
		switch t := tks[i]; t.T {
		case fract.Operator:
			if !opr {
				fract.IPanic(t, obj.SyntaxPanic, "Operator overflow!")
			}
			opr = false
			if b > 0 {
				part = append(part, t)
			} else {
				procs = append(procs, part)
				procs = append(procs, obj.Tokens{t})
				part = obj.Tokens{}
			}
		case fract.Value, fract.Name, fract.Comma, fract.Brace, fract.Loop, fract.In:
			if t.T == fract.Brace {
				switch t.V {
				case "(", "[", "{":
					b++
				default:
					b--
				}
			}
			if b == 0 && t.T == fract.Comma {
				fract.IPanic(t, obj.SyntaxPanic, "Invalid syntax!")
			}
			if i > 0 {
				if lt := tks[i-1]; (lt.T == fract.Name || lt.T == fract.Value) && (t.T == fract.Name || t.T == fract.Value) {
					fract.IPanic(t, obj.SyntaxPanic, "Invalid syntax!")
				}
			}
			part = append(part, t)
			opr = t.T != fract.Comma && (t.T != fract.Brace || t.T == fract.Brace && t.V != "[" && t.V != "(" && t.V != "{") && i < len(tks)-1
		default:
			fract.IPanic(t, obj.SyntaxPanic, "Invalid syntax!")
		}
	}
	if len(part) != 0 {
		procs = append(procs, part)
	}
	return procs
}

// decomposeBrace returns range tokens and index of first parentheses.
// Remove range tokens from original tokens.
func decomposeBrace(tks *obj.Tokens, ob, cb string) (obj.Tokens, int) {
	fst := -1
	for i, t := range *tks {
		if t.T == fract.Brace && t.V == ob {
			fst = i
			break
		}
	}
	// Skip find close parentheses and result ready steps
	// if open parentheses is not found.
	if fst == -1 {
		return nil, -1
	}
	// Find close parentheses.
	c, l := 1, 0
	for i := fst + 1; i < len(*tks); i++ {
		tk := (*tks)[i]
		if tk.T == fract.Brace {
			switch tk.V {
			case ob:
				c++
			case cb:
				c--
			}
			if c == 0 {
				break
			}
		}
		l++
	}
	rg := tks.Sub(fst+1, l)
	// Remove range from original tokens.
	tks.Rem(fst, (fst+l+1)-fst+1)
	if rg == nil {
		return nil, fst
	}
	return *rg, fst
}

// procIndex process array index by length.
func procIndex(len, i int) int {
	if i >= 0 {
		if i >= len {
			return -1
		}
		return i
	}
	i = len + i
	if i < 0 || i >= len {
		return -1
	}
	return i
}

// IsBlock returns true if tokens is block start, return false if not.
func IsBlock(tks obj.Tokens) bool {
	if tks[0].T == fract.Macro { // Remove macro token.
		tks = tks[1:]
	}
	switch tks[0].T {
	case fract.If, fract.Loop, fract.Func, fract.Try:
		return true
	case fract.Protected:
		if len(tks) > 1 && tks[1].T == fract.Func {
			return true
		}
	}
	return false
}

// nextopr find index of priority operator and returns index of operator
// if found, returns -1 if not.
func nextopr(tks []obj.Tokens) int {
	high, mid, low := -1, -1, -1
	for i, tslc := range tks {
		switch tslc[0].V {
		case "<<", ">>":
			return i
		case "**":
			return i
		case "%":
			return i
		case "*", "/", "\\", "//", "\\\\":
			if high == -1 {
				high = i
			}
		case "+", "-":
			if low == -1 {
				low = i
			}
		case "&", "|":
			if mid == -1 {
				mid = i
			}
		}
	}
	if high != -1 {
		return high
	} else if mid != -1 {
		return mid
	} else if low != -1 {
		return low
	}
	return -1
}

// findConditionOpr return next condition operator.
func findConditionOpr(tks obj.Tokens) (int, obj.Token) {
	bc := 0
	for i, t := range tks {
		if t.T == fract.Brace {
			switch t.V {
			case "{", "[", "(":
				bc++
			default:
				bc--
			}
		}
		if bc > 0 {
			continue
		}
		if t.T == fract.Operator {
			switch t.V {
			case "&&", "||", "==", "<>", ">", "<", "<=", ">=":
				return i, t
			}
		}
	}
	var tk obj.Token
	return -1, tk
}

// Find next or condition operator index and return if find, return -1 if not.
func nextConditionOpr(tks obj.Tokens, pos int, opr string) int {
	bc := 0
	for ; pos < len(tks); pos++ {
		t := tks[pos]
		if t.T == fract.Brace {
			switch t.V {
			case "{", "[", "(":
				bc++
			default:
				bc--
			}
		}
		if bc > 0 {
			continue
		}
		if t.T == fract.Operator && t.V == opr {
			return pos
		}
	}
	return -1
}

// conditionalProcesses returns conditional expressions by operators.
func conditionalProcesses(tks obj.Tokens, opr string) []obj.Tokens {
	var exps []obj.Tokens
	last := 0
	i := nextConditionOpr(tks, last, opr)
	for i != -1 {
		if i-last == 0 {
			fract.IPanic(tks[last], obj.SyntaxPanic, "Condition expression is cannot given!")
		}
		exps = append(exps, *tks.Sub(last, i-last))
		last = i + 1
		i = nextConditionOpr(tks, last, opr) // Find next.
		if i == len(tks)-1 {
			fract.IPanic(tks[len(tks)-1], obj.SyntaxPanic, "Operator overflow!")
		}
	}
	if last != len(tks) {
		exps = append(exps, *tks.Sub(last, len(tks)-last))
	}
	return exps
}

//! Built-in functions should have a lowercase names.

// ApplyBuildInFunctions to parser source.
func (p *Parser) AddBuiltInFuncs() {
	p.funcs = append(p.funcs,
		function{ // print function.
			name:              "print",
			protected:         true,
			tks:               nil,
			defaultParamCount: 2,
			params: []obj.Param{{
				Name: "value",
				Default: obj.Value{
					D: []obj.Data{
						{D: "", T: obj.VStr},
					},
				},
			}, {
				Name: "fin",
				Default: obj.Value{
					D: []obj.Data{
						{D: "\n", T: obj.VStr},
					},
				},
			}},
		}, function{ // input function.
			name:              "input",
			protected:         true,
			tks:               nil,
			defaultParamCount: 1,
			params: []obj.Param{{
				Name: "message",
				Default: obj.Value{
					D: []obj.Data{
						{D: "", T: obj.VStr},
					},
				},
			}},
		}, function{ // exit function.
			name:              "exit",
			protected:         true,
			tks:               nil,
			defaultParamCount: 1,
			params: []obj.Param{{
				Name: "code",
				Default: obj.Value{
					D: []obj.Data{{D: "0"}},
				},
			}},
		}, function{ // len function.
			name:              "len",
			protected:         true,
			tks:               nil,
			defaultParamCount: 0,
			params: []obj.Param{
				{Name: "object"},
			},
		}, function{ // range function.
			name:              "range",
			protected:         true,
			tks:               nil,
			defaultParamCount: 1,
			params: []obj.Param{
				{Name: "start"},
				{Name: "to"},
				{
					Name: "step",
					Default: obj.Value{
						D: []obj.Data{{D: "1", T: obj.VInt}},
					},
				},
			},
		}, function{ // calloc function.
			name:              "calloc",
			protected:         true,
			tks:               nil,
			defaultParamCount: 0,
			params: []obj.Param{
				{Name: "size"},
			},
		}, function{ // realloc function.
			name:              "realloc",
			protected:         true,
			tks:               nil,
			defaultParamCount: 0,
			params: []obj.Param{
				{Name: "base"},
				{Name: "size"},
			},
		}, function{ // memset function.
			name:              "memset",
			protected:         true,
			tks:               nil,
			defaultParamCount: 0,
			params: []obj.Param{
				{Name: "mem"},
				{Name: "val"},
			},
		}, function{ // string function.
			name:              "string",
			protected:         true,
			tks:               nil,
			defaultParamCount: 1,
			params: []obj.Param{
				{Name: "object"},
				{
					Name: "type",
					Default: obj.Value{
						D: []obj.Data{
							{D: "parse", T: obj.VStr},
						},
					},
				},
			},
		}, function{ // int function.
			name:              "int",
			protected:         true,
			tks:               nil,
			defaultParamCount: 1,
			params: []obj.Param{
				{Name: "object"},
				{
					Name: "type",
					Default: obj.Value{
						D: []obj.Data{
							{D: "parse", T: obj.VStr},
						},
					},
				},
			},
		}, function{ // float function.
			name:              "float",
			protected:         true,
			tks:               nil,
			defaultParamCount: 0,
			params: []obj.Param{
				{Name: "object"},
			},
		}, function{ // append function.
			name:              "append",
			protected:         true,
			tks:               nil,
			defaultParamCount: 0,
			params: []obj.Param{
				{Name: "dest"},
				{Name: "src", Params: true},
			},
		},
	)
}

// TODO: Add "match" keyword.
//! A change added here(especially added a code block) must also be compatible with "imports.go" and
//! add to "isBlock" function of parser.

// process tokens and returns true if block end, returns false if not and returns keyword state.
func (p *Parser) process(tks obj.Tokens) uint8 {
	//tks = append(obj.Tokens{}, tks...)
	switch fst := tks[0]; fst.T {
	case fract.Value, fract.Brace, fract.Name:
		if fst.T == fract.Name {
			bc := 0
			for _, t := range tks {
				if t.T == fract.Brace {
					switch t.V {
					case " {", "[", "(":
						bc++
					default:
						bc--
					}
				}
				if bc > 0 {
					continue
				}
				if t.T == fract.Operator {
					switch t.V {
					case "=", "+=", "-=", "*=", "/=", "%=", "^=", "<<=", ">>=", "|=", "&=":
						p.varset(tks)
						return fract.None
					}
				}
			}
		}
		// Print value if live interpreting.
		if v := p.procVal(tks); fract.InteractiveSh {
			if v.Print() {
				fmt.Println()
			}
		}
	case fract.Protected: // Protected declaration.
		if len(tks) < 2 {
			fract.IPanic(fst, obj.SyntaxPanic, "Define is not given!")
		}
		second := tks[1]
		tks = tks[1:]
		switch second.T {
		case fract.Var: // Variable definition.
			p.vardec(tks, true)
		case fract.Func: // Function definition.
			p.funcdec(tks, true)
		default:
			fract.IPanic(second, obj.SyntaxPanic, "Can protect only deletable objects!")
		}
	case fract.Var: // Variable definition.
		p.vardec(tks, false)
	case fract.Delete: // Delete from memory.
		p.procDel(tks)
	case fract.If: // if-elif-else.
		return p.procIf(tks)
	case fract.Loop: // Loop definition.
		p.loopCount++
		state := p.procLoop(tks)
		p.loopCount--
		return state
	case fract.Break: // Break loop.
		if p.loopCount < 1 {
			fract.IPanic(fst, obj.SyntaxPanic, "Break keyword only used in loops!")
		}
		return fract.LOOPBreak
	case fract.Continue: // Continue loop.
		if p.loopCount < 1 {
			fract.IPanic(fst, obj.SyntaxPanic, "Continue keyword only used in loops!")
		}
		return fract.LOOPContinue
	case fract.Ret: // Return.
		if p.funcCount < 1 {
			fract.IPanic(fst, obj.SyntaxPanic, "Return keyword only used in functions!")
		}
		if len(tks) > 1 {
			value := p.procVal(tks[1:])
			p.retVal = &value
		} else {
			p.retVal = nil
		}
		return fract.FUNCReturn
	case fract.Func: // Function definiton.
		p.funcdec(tks, false)
	case fract.Try: // Try-Catch.
		return p.procTryCatch(tks)
	case fract.Import: // Import.
		p.procImport(tks)
	case fract.Macro: // Macro.
		return p.procMacro(tks)
	case fract.Defer: // Defer.
		if l := len(tks); l < 2 {
			fract.IPanic(tks[0], obj.SyntaxPanic, "Function is not given!")
		} else if t := tks[l-1]; t.T != fract.Brace && t.V != ")" {
			fract.IPanicC(tks[0].F, tks[0].Ln, tks[0].Col+len(tks[0].V), obj.SyntaxPanic, "Invalid syntax!")
		}
		var vtks obj.Tokens
		bc := 0
		for i := len(tks) - 1; i >= 0; i-- {
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
			vtks = tks[1:i]
			break
		}
		if len(vtks) == 0 && bc == 0 {
			fract.IPanic(tks[1], obj.SyntaxPanic, "Invalid syntax!")
		}
		// Function call.
		v := p.procValPart(false, vtks)
		if v.Arr || v.D[0].T != obj.VFunc {
			fract.IPanic(tks[len(vtks)], obj.ValuePanic, "Value is not function!")
		}
		defers = append(defers, p.funcCallModel(v.D[0].D.(function), tks[len(vtks):]))
	default:
		fract.IPanic(fst, obj.SyntaxPanic, "Invalid syntax!")
	}
	return fract.None
}
