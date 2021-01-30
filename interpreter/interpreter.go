package interpreter

import (
	"fmt"

	"../lexer"
	"../objects"
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
	var preter *Interpreter = new(Interpreter)
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

	for !i.lexer.Finished {
		var tokens list.List = i.lexer.Next()
		for index := 0; index < tokens.Len(); index++ {
			fmt.Print(tokens.At(index).(objects.Token).Value)
		}
		fmt.Println()
	}
}
