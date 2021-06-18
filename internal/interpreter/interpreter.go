package interpreter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/fract-lang/fract/internal/lexer"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
)

// ReadyFile returns instance of source file by path.
func ReadyFile(filepath string) *objects.SourceFile {
	file, _ := os.Open(filepath)
	content, _ := os.ReadFile(filepath)
	return &objects.SourceFile{
		Lines: ReadyLines(strings.Split(string(content), "\n")),
		Path:  filepath,
		File:  file,
	}
}

// ReadyLines returns lines processed to lexing.
func ReadyLines(lines []string) []string {
	readyLines := make([]string, len(lines))
	for index, line := range lines {
		readyLines[index] = strings.TrimRight(line, " \t\n\r")
	}
	return readyLines
}

var (
	defers []functionCall
)

// Interpreter of Fract.
type Interpreter struct {
	variables         []objects.Variable
	functions         []objects.Function
	macroDefines      []objects.Variable
	funcTempVariables int // Count of function temporary variables.
	loopCount         int
	functionCount     int
	index             int
	returnValue       *objects.Value // Pointer of last return value.

	Lexer   *lexer.Lexer
	Tokens  [][]objects.Token // All Tokens of code file.
	Imports []importInfo
}

// New returns instance of interpreter related to file.
func New(path, fpath string) *Interpreter {
	return &Interpreter{
		Lexer: &lexer.Lexer{
			File: ReadyFile(fpath),
			Line: 1,
		},
	}
}

// NewStdin returns new instance of interpreter from standard input.
func NewStdin(path string) *Interpreter {
	return &Interpreter{
		Lexer: &lexer.Lexer{
			File: &objects.SourceFile{
				Path:  "<stdin>",
				File:  nil,
				Lines: nil,
			},
			Line: 1,
		},
	}
}

// ready interpreter to process.
func (i *Interpreter) ready() {
	/// Tokenize all lines.
	for !i.Lexer.Finished {
		if cacheTokens := i.Lexer.Next(); cacheTokens != nil {
			i.Tokens = append(i.Tokens, cacheTokens)
		}
	}
	// Change blocks.
	// TODO: Check "end" keyword alonity.
	count := 0
	macroBlockCount := 0
	last := -1
	for index, tokens := range i.Tokens {
		if first := tokens[0]; first.Type == fract.TypeBlockEnd {
			count--
			if count < 0 {
				fract.Error(first, "The extra block end defined!")
			}
		} else if first.Type == fract.TypeMacro {
			if parser.IsBlockStatement(tokens) {
				macroBlockCount++
				if macroBlockCount == 1 {
					last = index
				}
			} else if tokens[1].Type == fract.TypeBlockEnd {
				macroBlockCount--
				if macroBlockCount < 0 {
					fract.Error(first, "The extra block end defined!")
				}
			}
		} else if parser.IsBlockStatement(tokens) {
			count++
			if count == 1 {
				last = index
			}
		}
	}
	if count > 0 || macroBlockCount > 0 { // Check blocks.
		fract.Error(i.Tokens[last][0], "Block is expected ending...")
	}
}

func (i *Interpreter) Interpret() {
	if i.Lexer.File.Path == "<stdin>" {
		// Interpret all lines.
		for i.index = 0; i.index < len(i.Tokens); i.index++ {
			i.processTokens(i.Tokens[i.index])
			runtime.GC()
		}
		goto final
	}
	// Lexer is finished.
	if i.Lexer.Finished {
		return
	}

	i.ready()
	{
		//* Import local directory.
		dir, _ := os.Getwd()
		if pdir := path.Dir(i.Lexer.File.Path); pdir != "." {
			dir = path.Join(dir, pdir)
		}
		content, err := ioutil.ReadDir(dir)
		if err == nil {
			_, mainName := filepath.Split(i.Lexer.File.Path)
			for _, current := range content {
				// Skip directories.
				if current.IsDir() || !strings.HasSuffix(current.Name(), fract.FractExtension) || current.Name() == mainName {
					continue
				}
				source := New(dir, path.Join(dir, current.Name()))
				source.ApplyEmbedFunctions()
				source.Import()
				i.functions = append(i.functions, source.functions...)
				i.variables = append(i.variables, source.variables...)
				i.macroDefines = append(i.macroDefines, source.macroDefines...)
			}
		}
	}

	// Interpret all lines.
	for i.index = 0; i.index < len(i.Tokens); i.index++ {
		i.processTokens(i.Tokens[i.index])
		runtime.GC()
	}

final:
	for index := len(defers) - 1; index >= 0; index-- {
		defers[index].call()
	}
}

// skipBlock skip to block end.
func (i *Interpreter) skipBlock(ifBlock bool) {
	if ifBlock {
		if parser.IsBlockStatement(i.Tokens[i.index]) {
			i.index++
		} else {
			return
		}
	}
	count := 1
	i.index--
	for {
		i.index++
		tokens := i.Tokens[i.index]
		if first := tokens[0]; first.Type == fract.TypeBlockEnd {
			count--
		} else if first.Type == fract.TypeMacro {
			if parser.IsBlockStatement(tokens) {
				count++
			} else if tokens[1].Type == fract.TypeBlockEnd {
				count--
			}
		} else if parser.IsBlockStatement(tokens) {
			count++
		}
		if count == 0 {
			return
		}
	}
}

// TYPES
// 'f' -> Function.
// 'v' -> Variable.
func (i *Interpreter) defineByName(name objects.Token) (int, rune, *Interpreter) {
	index, source := i.functionIndexByName(name)
	if index != -1 {
		return index, 'f', source
	}
	index, source = i.variableIndexByName(name)
	if index != -1 {
		return index, 'v', source
	}
	return -1, '-', nil
}

func (i *Interpreter) definedName(name objects.Token) int {
	if name.Value[0] == '-' { // Ignore minus.
		name.Value = name.Value[1:]
	}
	for _, current := range i.functions {
		if current.Name == name.Value {
			return current.Line
		}
	}
	for _, current := range i.variables {
		if current.Name == name.Value {
			return current.Line
		}
	}
	return -1
}

//! This code block very like to variableIndexByName function. If you change here, probably you must change there too.

// functionIndexByName returns index of function by name.
func (i *Interpreter) functionIndexByName(name objects.Token) (int, *Interpreter) {
	if name.Value[0] == '-' { // Ignore minus.
		name.Value = name.Value[1:]
	}
	if index := strings.Index(name.Value, "."); index != -1 {
		if i.importIndexByName(name.Value[:index]) == -1 {
			fract.Error(name, "'"+name.Value[:index]+"' is not defined!")
		}
		i = i.Imports[i.importIndexByName(name.Value[:index])].Source
		name.Value = name.Value[index+1:]
		for index, current := range i.functions {
			if (current.Tokens == nil || unicode.IsUpper(rune(current.Name[0]))) && current.Name == name.Value {
				return index, i
			}
		}
		return -1, nil
	}
	for index, current := range i.functions {
		if current.Name == name.Value {
			return index, i
		}
	}
	return -1, nil
}

//! This code block very like to functionIndexByName function. If you change here, probably you must change there too.

// variableIndexByName returns index of variable by name.
func (i *Interpreter) variableIndexByName(name objects.Token) (int, *Interpreter) {
	if name.Value[0] == '-' { // Ignore minus.
		name.Value = name.Value[1:]
	}
	if index := strings.Index(name.Value, "."); index != -1 {
		if iindex := i.importIndexByName(name.Value[:index]); iindex == -1 {
			fract.Error(name, "'"+name.Value[:index]+"' is not defined!")
		} else {
			i = i.Imports[iindex].Source
		}
		name.Value = name.Value[index+1:]
		for index, current := range i.variables {
			if (current.Line == -1 || unicode.IsUpper(rune(current.Name[0]))) && current.Name == name.Value {
				return index, i
			}
		}
		return -1, nil
	}
	for index, current := range i.variables {
		if current.Name == name.Value {
			return index, i
		}
	}
	return -1, nil
}

// importIndexByName returns index of import by name.
func (i *Interpreter) importIndexByName(name string) int {
	for index, current := range i.Imports {
		if current.Name == name {
			return index
		}
	}
	return -1
}

//! Embed functions should have a lowercase names.
// TODO: Add append function.
// TODO: Add copy function.

// ApplyEmbedFunctions to interpreter source.
func (i *Interpreter) ApplyEmbedFunctions() {
	i.functions = append(i.functions,
		objects.Function{ // print function.
			Name:                  "print",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 2,
			Parameters: []objects.Parameter{
				{
					Name: "value",
					Default: objects.Value{
						Content: []objects.Data{
							{
								Type: objects.VALString,
							},
						},
					},
				},
				{
					Name: "fin",
					Default: objects.Value{
						Content: []objects.Data{
							{
								Data: "\n",
								Type: objects.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // input function.
			Name:                  "input",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []objects.Parameter{
				{
					Name: "message",
					Default: objects.Value{
						Content: []objects.Data{
							{
								Data: "",
								Type: objects.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // exit function.
			Name:                  "exit",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []objects.Parameter{
				{
					Name: "code",
					Default: objects.Value{
						Content: []objects.Data{{Data: "0"}},
					},
				},
			},
		},
		objects.Function{ // len function.
			Name:                  "len",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: []objects.Parameter{
				{
					Name: "object",
				},
			},
		},
		objects.Function{ // range function.
			Name:                  "range",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []objects.Parameter{
				{
					Name: "start",
				},
				{
					Name: "to",
				},
				{
					Name: "step",
					Default: objects.Value{
						Content: []objects.Data{{Data: "1"}},
					},
				},
			},
		},
		objects.Function{ // make function.
			Name:                  "make",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: []objects.Parameter{
				{
					Name: "size",
				},
			},
		},
		objects.Function{ // string function.
			Name:                  "string",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []objects.Parameter{
				{
					Name: "object",
				},
				{
					Name: "type",
					Default: objects.Value{
						Content: []objects.Data{
							{
								Data: "parse",
								Type: objects.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // int function.
			Name:                  "int",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 1,
			Parameters: []objects.Parameter{
				{
					Name: "object",
				},
				{
					Name: "type",
					Default: objects.Value{
						Content: []objects.Data{
							{
								Data: "parse",
								Type: objects.VALString,
							},
						},
					},
				},
			},
		},
		objects.Function{ // float function.
			Name:                  "float",
			Protected:             true,
			Tokens:                nil,
			DefaultParameterCount: 0,
			Parameters: []objects.Parameter{
				{
					Name: "object",
				},
			},
		},
	)
}

// TODO: Add "match" keyword.
//! A change added here(especially added a code block) must also be compatible with "import.go" and
//! add to "isBlockStatement.go" of parser.

// processTokens returns true if block end, returns false if not and returns keyword state.
func (i *Interpreter) processTokens(tokens []objects.Token) uint8 {
	tokens = append([]objects.Token{}, tokens...)
	switch first := tokens[0]; first.Type {
	case
		fract.TypeValue,
		fract.TypeBrace,
		fract.TypeName:
		if first.Type == fract.TypeName {
			brace := 0
			for _, current := range tokens {
				if current.Type == fract.TypeBrace {
					if current.Value == "{" || current.Value == "[" || current.Value == "(" {
						brace++
					} else {
						brace--
					}
				}
				if brace > 0 {
					continue
				}
				if current.Type == fract.TypeOperator &&
					(current.Value == "=" ||
						current.Value == grammar.AdditionAssignment ||
						current.Value == grammar.SubtractionAssignment ||
						current.Value == grammar.MultiplicationAssignment ||
						current.Value == grammar.DivisionAssignment ||
						current.Value == grammar.ModulusAssignment ||
						current.Value == grammar.XOrAssignment ||
						current.Value == grammar.LeftBinaryShiftAssignment ||
						current.Value == grammar.RightBinaryShiftAssignment ||
						current.Value == grammar.InclusiveOrAssignment ||
						current.Value == grammar.AndAssignment) { // Variable setting.
					i.processVariableSet(tokens)
					return fract.TypeNone
				}
			}
		}
		// Print value if live interpreting.
		if value := i.processValue(tokens); fract.InteractiveShell {
			if value.Print() {
				fmt.Println()
			}
		}
	case fract.TypeProtected: // Protected declaration.
		if len(tokens) < 2 {
			fract.Error(first, "Protected but what is it protected?")
		}
		second := tokens[1]
		tokens = tokens[1:]
		if second.Type == fract.TypeVariable { // Variable definition.
			i.processVariableDeclaration(tokens, true)
		} else if second.Type == fract.TypeFunction { // Function definition.
			i.processFunctionDeclaration(tokens, true)
		} else {
			fract.Error(second, "Syntax error, you can protect only deletable objects!")
		}
	case fract.TypeVariable: // Variable definition.
		i.processVariableDeclaration(tokens, false)
	case fract.TypeDelete: // Delete from memory.
		i.processDelete(tokens)
	case fract.TypeIf: // if-elif-else.
		return i.processIf(tokens)
	case fract.TypeLoop: // Loop definition.
		i.loopCount++
		state := i.processLoop(tokens)
		i.loopCount--
		return state
	case fract.TypeBreak: // Break loop.
		if i.loopCount == 0 {
			fract.Error(first, "Break keyword only used in loops!")
		}
		return fract.LOOPBreak
	case fract.TypeContinue: // Continue loop.
		if i.loopCount == 0 {
			fract.Error(first, "Continue keyword only used in loops!")
		}
		return fract.LOOPContinue
	case fract.TypeReturn: // Return.
		if i.functionCount == 0 {
			fract.Error(first, "Return keyword only used in functions!")
		}
		if len(tokens) > 1 {
			value := i.processValue(tokens[1:])
			i.returnValue = &value
		} else {
			i.returnValue = nil
		}
		return fract.FUNCReturn
	case fract.TypeFunction: // Function definiton.
		i.processFunctionDeclaration(tokens, false)
	case fract.TypeTry: // Try-Catch.
		return i.processTryCatch(tokens)
	case fract.TypeImport: // Import.
		i.processImport(tokens)
	case fract.TypeMacro: // Macro.
		return i.processMacro(tokens)
	case fract.TypeDefer: // Defer.
		if l := len(tokens); l < 2 {
			fract.Error(tokens[0], "Function is not defined!")
		} else if tokens[1].Type != fract.TypeName {
			fract.Error(tokens[1], "Invalid syntax!")
		} else if l < 3 {
			fract.Error(tokens[1], "Invalid syntax!")
		} else if tokens[2].Type != fract.TypeBrace || tokens[2].Value != "(" {
			fract.Error(tokens[2], "Invalid syntax!")
		}
		defers = append(defers, i.processFunctionCallModel(tokens[1:]))
	default:
		fract.Error(first, "Invalid syntax!")
	}
	return fract.TypeNone
}
