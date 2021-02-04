/*
	GENERIC FUNCTIONS
*/

package interpreter

import (
	"strings"

	"../lexer"
	"../objects"
	"../utilities/fs"
	"../utilities/vector"
)

// New Create new instance of Parser.
// path Path of destination file.
// type Type of file.
func New(path string, _type int) *Interpreter {
	preter := new(Interpreter)
	preter.lexer = lexer.New(ReadyFile(path))
	preter.vars = vector.New()
	preter.Type = _type
	return preter
}

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
		text := strings.TrimSpace(lines[index])
		if text == "" {
			continue
		}
		readyLines.Append(objects.CodeLine{Line: index + 1, Text: text})
	}
	return readyLines
}
