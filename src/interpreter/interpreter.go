package interpreter

import (
	"fmt"
	"strings"

	"../fract"
	"../fract/arithmetic"
	"../grammar"
	"../lexer"
	"../objects"
	"../parser"
	"../utilities/fs"
	"../utilities/list"
)

// Interpreter Interprater of Fract.
type Interpreter struct {
	/* PRIVATE */

	// Parser of this file.
	lexer *lexer.Lexer

	/* PUBLIC */

	/* Type of file. */
	Type int
}

// *********************
//       PRIVATE
// *********************

// processValue Process value.
// tokens Tokens.
func (i *Interpreter) processValue(tokens *list.List) objects.Value {
	var (
		value  objects.Value
		avalue float64 = 0
	)
	value.Content = ""
	value.Type = fract.VTInteger

	// Decompose arithmetic operations
	operations := parser.DecomposeArithmeticProcesses(tokens)
	for index := 0; index < operations.Len(); index++ {
		operation := operations.Vals[index].(objects.ArithmeticProcess)

		/* Check values. */
		if !arithmetic.IsFloat(operation.First.Value) {
			fract.Error(operation.First,
				"This is not a arithmetic value!: "+operation.First.Value)
		} else if !arithmetic.IsFloat(operation.Second.Value) {
			fract.Error(operation.Second,
				"This is not a arithmetic value!: "+operation.Second.Value)
		}

		if strings.Index(operation.First.Value, grammar.TokenDot) != -1 ||
			strings.Index(operation.Second.Value, grammar.TokenDot) != -1 {
			value.Type = fract.VTFloat
		}

		avalue += arithmetic.SolveArithmeticProcess(operation)
	}
	value.Content = arithmetic.FloatToString(avalue)

	return value
}

// *********************
//        PUBLIC
// *********************

// ReadyFile Create instance of code file.
// path Path of file.
func ReadyFile(path string) objects.CodeFile {
	var file objects.CodeFile
	file.Lines = ReadyLines(fs.ReadAllLines(path))
	file.Path = path
	file.File = fs.OpenFile(path)
	return file
}

// ReadyLines Ready lines to process.
// lines Lines to ready.
func ReadyLines(lines []string) []objects.CodeLine {
	var readyLines []objects.CodeLine
	for index := 0; index < len(lines); index++ {
		readyLines = append(readyLines, objects.CodeLine{Line: index + 1, Text: lines[index]})
	}
	return readyLines
}

// New Create new instance of Parser.
// path Path of destination file.
// type Type of file.
func New(path string, _type int) *Interpreter {
	preter := new(Interpreter)
	preter.lexer = lexer.New(ReadyFile(path))
	preter.Type = _type
	return preter
}

// Interpret Interpret file.
func (i *Interpreter) Interpret() {
	// Lexer is finished.
	if i.lexer.Finished {
		return
	}

	/* Interpret all lines. */
	for !i.lexer.Finished {
		tokens := i.lexer.Next()
		first := tokens.Vals[0].(objects.Token)

		if first.Type == fract.TypeValue {
			fmt.Println(i.processValue(&tokens).Content)
		} else {
			fract.Error(first, "What is this?: "+first.Value)
		}
	}
}
