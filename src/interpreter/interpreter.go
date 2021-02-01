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
	"../utilities/vector"
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
func (i *Interpreter) processValue(tokens *vector.Vector) objects.Value {
	/* Check parentheses range. */
	for true {
		_range, found := parser.DecomposeBrace(tokens)

		/* Parentheses are not found! */
		if found == -1 || len(_range.Vals) == 0 {
			break
		}

		var _token objects.Token
		_token.Value = i.processValue(&_range).Content
		_token.Type = fract.TypeValue
		tokens.Insert(found, _token)
	}

	var (
		value     objects.Value
		operation objects.ArithmeticProcess
		avalue    float64 = 0
	)
	value.Content = ""
	value.Type = fract.VTInteger

	// Decompose arithmetic operations
	operations := parser.DecomposeArithmeticProcesses(tokens)
	for index := 0; index < len(operations.Vals); index++ {
		_token := operations.Vals[index].(objects.Token)
		if operation.First.Value == "" {
			operation.First = _token
			continue
		} else if operation.Operator.Value == "" {
			operation.Operator = _token
			continue
		}
		operation.Second = _token
		avalue = arithmetic.SolveArithmeticProcess(operation)

		/* Reset to defaults. */
		operation.First = operation.Second
		operation.First.Value = arithmetic.FloatToString(avalue)
		operation.Operator.Value = ""
		operation.Second.Value = ""
	}
	// If only one value.
	if operations.Len() == 1 {
		avalue, _ = arithmetic.ToDouble(operations.First().(objects.Token).Value)
	}

	// Set value.
	value.Content = arithmetic.FloatToString(avalue)

	/* Set type to float if... */
	if value.Type != fract.VTFloat &&
		(strings.Index(value.Content, grammar.TokenDot) != -1 ||
			strings.Index(value.Content, grammar.TokenDot) != -1) {
		value.Type = fract.VTFloat
	}

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
			if first.Type == fract.TypeBrace {
				fract.Error(first, "Statement are don't starts with brackets!")
			}
			fract.Error(first, "What is this?: "+first.Value)
		}
	}
}
