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
		if found == -1 {
			break
		}

		var _token objects.Token
		_token.Value = i.processValue(&_range).Content
		_token.Type = fract.TypeValue
		tokens.Insert(found, _token)
	}

	var value objects.Value
	value.Content = ""
	value.Type = fract.VTInteger

	// Decompose arithmetic operations.
	operations := parser.DecomposeArithmeticProcesses(tokens)

	// Process arithmetic operation.
	priorityIndex := parser.IndexProcessPriority(&operations)
	for priorityIndex != -1 {
		var operation objects.ArithmeticProcess
		operation.First = operations.Vals[priorityIndex-1].(objects.Token)
		operation.Operator = operations.Vals[priorityIndex].(objects.Token)
		operation.Second = operations.Vals[priorityIndex+1].(objects.Token)

		_token := operations.Vals[priorityIndex-1].(objects.Token)
		operations.RemoveRange(priorityIndex-1, 3)
		_type, result := arithmetic.SolveArithmeticProcess(operation)
		value.Type = _type
		_token.Value = arithmetic.TypeToString(_type, result)
		operations.Insert(priorityIndex-1, _token)

		// Find next operator.
		priorityIndex = parser.IndexProcessPriority(&operations)
	}

	// Set value.
	_value, _ := arithmetic.ToDouble(operations.First().(objects.Token).Value)
	value.Content = arithmetic.TypeToString(value.Type, _value)

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
func ReadyLines(lines []string) *vector.Vector {
	readyLines := vector.New()
	for index := 0; index < len(lines); index++ {
		readyLines.Append(objects.CodeLine{Line: index + 1, Text: lines[index]})
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

		// Skip this loop if tokens are empty.
		if tokens.Len() == 0 {
			continue
		}

		first := tokens.Vals[0].(objects.Token)

		if first.Type == fract.TypeValue || first.Type == fract.TypeBrace {
			fmt.Println(i.processValue(&tokens).Content)
		} else {
			if first.Type == fract.TypeBrace {
				fract.Error(first, "Statement are don't starts with brackets!")
			}
			fract.Error(first, "What is this?: "+first.Value)
		}
	}
}
