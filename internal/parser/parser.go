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

// File returns instance of source file by path.
func File(fp string) *obj.File {
	f, _ := os.Open(fp)
	bytes, _ := os.ReadFile(fp)
	return &obj.File{
		Lns: Lines(strings.Split(string(bytes), "\n")),
		P:   fp,
		F:   f,
	}
}

// Lines returns ready lines processed to lexing.
func Lines(lns []string) []string {
	for i, ln := range lns {
		lns[i] = strings.TrimRightFunc(ln, func(r rune) bool { return r == '\r' })
	}
	return lns
}

var (
	defers []funcCall
)

// Parser of Fract.
type Parser struct {
	vars         []obj.Var
	funcs        []obj.Func
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
	return &Parser{
		funcTempVars: -1,
		L: &lex.Lex{
			F:  File(fp),
			Ln: 1,
		},
	}
}

// NewStdin returns new instance of parser from standard input.
func NewStdin() *Parser {
	return &Parser{
		funcTempVars: -1,
		L: &lex.Lex{
			F: &obj.File{
				P:   "<stdin>",
				F:   nil,
				Lns: nil,
			},
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
				fract.Error(tks[1], "Invalid syntax!")
			}
			bc--
			if bc < 0 {
				fract.Error(fst, "The extra block end defined!")
			}
		} else if fst.T == fract.Macro {
			if IsBlock(tks) {
				mbc++
				if mbc == 1 {
					lst = i
				}
			} else if tks[1].T == fract.End {
				if len(tks) > 2 {
					fract.Error(tks[2], "Invalid syntax!")
				}
				mbc--
				if mbc < 0 {
					fract.Error(fst, "The extra block end defined!")
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
		fract.Error(p.Tks[lst][0], "Block is expected ending...")
	}
}

func (p *Parser) Interpret() {
	if p.L.F.P == "<stdin>" {
		// Interpret all lines.
		for p.i = 0; p.i < len(p.Tks); p.i++ {
			p.process(p.Tks[p.i])
			runtime.GC()
		}
		goto end
	}
	// Lexer is finished.
	if p.L.Fin {
		return
	}
	p.ready()
	{
		//* Import local directory.
		dir, _ := os.Getwd()
		if pdir := path.Dir(p.L.F.P); pdir != "." {
			dir = path.Join(dir, pdir)
		}
		infos, err := ioutil.ReadDir(dir)
		if err == nil {
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
	}
	// Interpret all lines.
	for p.i = 0; p.i < len(p.Tks); p.i++ {
		p.process(p.Tks[p.i])
		runtime.GC()
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
				fract.Error(tk, "Only integer values can used in index access!")
			}
			pos, err := strconv.Atoi(d.String())
			if err != nil {
				fract.Error(tk, "Value out of range!")
			}
			if arr.Arr {
				pos = procIndex(len(arr.D), pos)
			} else {
				pos = procIndex(len(arr.D[0].String()), pos)
			}
			if pos == -1 {
				fract.Error(tk, "Index is out of range!")
			}
			i = append(i, pos)
		}
		return i
	}
	if val.D[0].T != obj.VInt {
		fract.Error(tk, "Only integer values can used in index access!")
	}
	pos, err := strconv.Atoi(val.String())
	if err != nil {
		fract.Error(tk, "Value out of range!")
	}
	if arr.Arr {
		pos = procIndex(len(arr.D), pos)
	} else {
		pos = procIndex(len(arr.D[0].String()), pos)
	}
	if pos == -1 {
		fract.Error(tk, "Index is out of range!")
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
	if name.Val[0] == '-' { // Ignore minus.
		name.Val = name.Val[1:]
	}
	for _, f := range p.funcs {
		if f.Name == name.Val {
			return f.Ln
		}
	}
	for _, v := range p.vars {
		if v.Name == name.Val {
			return v.Ln
		}
	}
	return -1
}

//! This code block very like to variableIndexByName function. If you change here, probably you must change there too.

// funcIndexByName returns index of function by name.
func (p *Parser) funcIndexByName(name obj.Token) (int, *Parser) {
	if name.Val[0] == '-' { // Ignore minus.
		name.Val = name.Val[1:]
	}
	if i := strings.Index(name.Val, "."); i != -1 {
		if p.importIndexByName(name.Val[:i]) == -1 {
			fract.Error(name, "'"+name.Val[:i]+"' is not defined!")
		}
		p = p.Imports[p.importIndexByName(name.Val[:i])].Src
		name.Val = name.Val[i+1:]
		for i, current := range p.funcs {
			if (current.Tks == nil || unicode.IsUpper(rune(current.Name[0]))) && current.Name == name.Val {
				return i, p
			}
		}
		return -1, nil
	}
	for j, f := range p.funcs {
		if f.Name == name.Val {
			return j, p
		}
	}
	return -1, nil
}

//! This code block very like to functionIndexByName function. If you change here, probably you must change there too.

// varIndexByName returns index of variable by name.
func (p *Parser) varIndexByName(name obj.Token) (int, *Parser) {
	if name.Val[0] == '-' { // Ignore minus.
		name.Val = name.Val[1:]
	}
	if i := strings.Index(name.Val, "."); i != -1 {
		if iindex := p.importIndexByName(name.Val[:i]); iindex == -1 {
			fract.Error(name, "'"+name.Val[:i]+"' is not defined!")
		} else {
			p = p.Imports[iindex].Src
		}
		name.Val = name.Val[i+1:]
		for i, v := range p.vars {
			if (v.Ln == -1 || unicode.IsUpper(rune(v.Name[0]))) && v.Name == name.Val {
				return i, p
			}
		}
		return -1, nil
	}
	for j, v := range p.vars {
		if v.Name == name.Val {
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
func checkArithmeticProcesses(tks []obj.Token) {
	var (
		opr bool
		b   int
	)
	for i := 0; i < len(tks); i++ {
		switch t := tks[i]; t.T {
		case fract.Operator:
			if !opr {
				fract.Error(t, "Operator spam!")
			}
			opr = false
		case fract.Value, fract.Name, fract.Comma, fract.Brace:
			switch t.T {
			case fract.Brace:
				if t.Val == "(" || t.Val == "[" || t.Val == "{" {
					b++
				} else {
					b--
				}
			case fract.Comma:
				if b == 0 {
					fract.Error(t, "Invalid syntax!")
				}
			}
			opr = i < len(tks)-1
		default:
			fract.Error(t, "Invalid syntax!")
		}
	}
	if tks[len(tks)-1].T == fract.Operator {
		fract.Error(tks[len(tks)-1], "Operator overflow!")
	}
}

// decomposeBrace returns range tokens and index of first parentheses.
// Remove range tokens from original tokens.
func decomposeBrace(tks *obj.Tokens, ob, cb string, noChk bool) ([]obj.Token, int) {
	fst := -1
	/* Find open parentheses. */
	if noChk {
		n := false
		for i, t := range *tks {
			if t.T == fract.Name {
				n = true
			} else if !n && t.T == fract.Brace && t.Val == ob {
				fst = i
				break
			} else {
				n = false
			}
		}
	} else {
		for i, t := range *tks {
			if t.T == fract.Brace && t.Val == ob {
				fst = i
				break
			}
		}
	}
	// Skip find close parentheses and result ready steps
	// if open parentheses is not found.
	if fst == -1 {
		return nil, -1
	}
	// Find close parentheses.
	c := 1
	l := 0
	for i := fst + 1; i < len(*tks); i++ {
		tk := (*tks)[i]
		if tk.T == fract.Brace {
			if tk.Val == ob {
				c++
			} else if tk.Val == cb {
				c--
			}
			if c == 0 {
				break
			}
		}
		l++
	}
	rg := tks.Sub(fst+1, l)
	// Bracket content is empty?
	if noChk && rg == nil {
		fract.Error((*tks)[fst], "Brackets content are empty!")
	}
	/* Remove range from original tokens. */
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
func IsBlock(tks []obj.Token) bool {
	if tks[0].T == fract.Macro { // Remove macro token.
		tks = tks[1:]
	}
	switch tks[0].T {
	case fract.If,
		fract.Loop,
		fract.Func,
		fract.Try:
		return true
	case fract.Protected:
		if len(tks) > 1 {
			if tks[1].T == fract.Func {
				return true
			}
		}
	}
	return false
}

// nextopr find index of priority operator and returns index of operator if found, returns -1 if not.
func nextopr(tks []obj.Token) int {
	bc := 0
	high := -1
	mid := -1
	low := -1
	for i, t := range tks {
		if t.T == fract.Brace {
			if t.Val == "[" || t.Val == "{" || t.Val == "(" {
				bc++
			} else {
				bc--
			}
		}
		if bc > 0 {
			continue
		}
		// Exponentiation or shifts.
		if t.Val == "<<" || t.Val == ">>" || t.Val == "**" {
			return i
		}
		switch t.Val {
		case "%": // Modulus.
			return i
		case "*", "/", "\\", "//", "\\\\": // Multiply or division.
			if high == -1 {
				high = i
			}
		case "+", "-": // Addition or subtraction.
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
func findConditionOpr(tks []obj.Token) (int, obj.Token) {
	bc := 0
	for i, t := range tks {
		if t.T == fract.Brace {
			if t.Val == "{" || t.Val == "[" || t.Val == "(" {
				bc++
			} else {
				bc--
			}
		} else if bc == 0 &&
			(t.T == fract.Operator && (t.Val == "&&" || t.Val == "||" ||
				t.Val == "==" || t.Val == "<>" || t.Val == ">" || t.Val == "<" ||
				t.Val == ">=" || t.Val == "<=")) || t.T == fract.In {
			return i, t
		}
	}
	var tk obj.Token
	return -1, tk
}

// Find next or condition operator index and return if find, return -1 if not.
func nextConditionOpr(tks []obj.Token, pos int, opr string) int {
	bc := 0
	for ; pos < len(tks); pos++ {
		t := tks[pos]
		if t.T == fract.Brace {
			if t.Val == "{" || t.Val == "[" || t.Val == "(" {
				bc++
			} else {
				bc--
			}
		}
		if bc > 0 {
			continue
		}
		if t.T == fract.Operator && t.Val == opr {
			return pos
		}
	}
	return -1
}

// conditionalProcesses returns conditional expressions by operators.
func conditionalProcesses(tks obj.Tokens, opr string) *[]obj.Tokens {
	var exps []obj.Tokens
	last := 0
	i := nextConditionOpr(tks, last, opr)
	for i != -1 {
		if i-last == 0 {
			fract.Error(tks[last], "Where is the condition?")
		}
		exps = append(exps, *tks.Sub(last, i-last))
		last = i + 1
		i = nextConditionOpr(tks, last, opr) // Find next.
		if i == len(tks)-1 {
			fract.Error(tks[len(tks)-1], "Operator defined, but for what?")
		}
	}
	if last != len(tks) {
		exps = append(exps, *tks.Sub(last, len(tks)-last))
	}
	return &exps
}

//! Embed functions should have a lowercase names.
// ApplyEmbedFunctions to parser source.
func (p *Parser) ApplyEmbedFunctions() {
	p.funcs = append(p.funcs,
		obj.Func{ // print function.
			Name:              "print",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 2,
			Params: []obj.Param{
				{
					Name: "value",
					Default: obj.Value{
						D: []obj.Data{
							{
								T: obj.VStr,
							},
						},
					},
				},
				{
					Name: "fin",
					Default: obj.Value{
						D: []obj.Data{
							{
								D: "\n",
								T: obj.VStr,
							},
						},
					},
				},
			},
		},
		obj.Func{ // input function.
			Name:              "input",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 1,
			Params: []obj.Param{
				{
					Name: "message",
					Default: obj.Value{
						D: []obj.Data{
							{
								D: "",
								T: obj.VStr,
							},
						},
					},
				},
			},
		},
		obj.Func{ // exit function.
			Name:              "exit",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 1,
			Params: []obj.Param{
				{
					Name: "code",
					Default: obj.Value{
						D: []obj.Data{{D: "0"}},
					},
				},
			},
		},
		obj.Func{ // len function.
			Name:              "len",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 0,
			Params: []obj.Param{
				{Name: "object"},
			},
		},
		obj.Func{ // range function.
			Name:              "range",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 1,
			Params: []obj.Param{
				{Name: "start"},
				{Name: "to"},
				{
					Name: "step",
					Default: obj.Value{
						D: []obj.Data{{D: "1"}},
					},
				},
			},
		},
		obj.Func{ // make function.
			Name:              "make",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 0,
			Params: []obj.Param{
				{Name: "size"},
			},
		},
		obj.Func{ // string function.
			Name:              "string",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 1,
			Params: []obj.Param{
				{Name: "object"},
				{
					Name: "type",
					Default: obj.Value{
						D: []obj.Data{
							{
								D: "parse",
								T: obj.VStr,
							},
						},
					},
				},
			},
		},
		obj.Func{ // int function.
			Name:              "int",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 1,
			Params: []obj.Param{
				{Name: "object"},
				{
					Name: "type",
					Default: obj.Value{
						D: []obj.Data{
							{
								D: "parse",
								T: obj.VStr,
							},
						},
					},
				},
			},
		},
		obj.Func{ // float function.
			Name:              "float",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 0,
			Params: []obj.Param{
				{Name: "object"},
			},
		},
		obj.Func{ // append function.
			Name:              "append",
			Protected:         true,
			Tks:               nil,
			DefaultParamCount: 0,
			Params: []obj.Param{
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
func (p *Parser) process(tks []obj.Token) uint8 {
	tks = append([]obj.Token{}, tks...)
	switch fst := tks[0]; fst.T {
	case
		fract.Value,
		fract.Brace,
		fract.Name:
		if fst.T == fract.Name {
			bc := 0
			for _, t := range tks {
				if t.T == fract.Brace {
					if t.Val == "{" || t.Val == "[" || t.Val == "(" {
						bc++
					} else {
						bc--
					}
				}
				if bc > 0 {
					continue
				}
				if t.T == fract.Operator &&
					(t.Val == "=" || t.Val == "+=" || t.Val == "-=" || t.Val == "*=" || t.Val == "/=" || t.Val == "%=" ||
						t.Val == "^=" || t.Val == "<<=" || t.Val == ">>=" || t.Val == "|=" || t.Val == "&=") { // Variable setting.
					p.varset(tks)
					return fract.None
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
			fract.Error(fst, "Protected but what is it protected?")
		}
		second := tks[1]
		tks = tks[1:]
		if second.T == fract.Var { // Variable definition.
			p.vardec(tks, true)
		} else if second.T == fract.Func { // Function definition.
			p.funcdec(tks, true)
		} else {
			fract.Error(second, "Syntax error, you can protect only deletable objects!")
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
		if p.loopCount == 0 {
			fract.Error(fst, "Break keyword only used in loops!")
		}
		return fract.LOOPBreak
	case fract.Continue: // Continue loop.
		if p.loopCount == 0 {
			fract.Error(fst, "Continue keyword only used in loops!")
		}
		return fract.LOOPContinue
	case fract.Ret: // Return.
		if p.funcCount == 0 {
			fract.Error(fst, "Return keyword only used in functions!")
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
			fract.Error(tks[0], "Function is not defined!")
		} else if tks[1].T != fract.Name {
			fract.Error(tks[1], "Invalid syntax!")
		} else if l < 3 {
			fract.Error(tks[1], "Invalid syntax!")
		} else if tks[2].T != fract.Brace || tks[2].Val != "(" {
			fract.Error(tks[2], "Invalid syntax!")
		}
		defers = append(defers, p.funcCallModel(tks[1:]))
	default:
		fract.Error(fst, "Invalid syntax!")
	}
	return fract.None
}
