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
	preter.tokens = vector.New()
	return preter
}

// ReadyFile Create instance of code file.
// path Path of file.
func ReadyFile(path string) objects.CodeFile {
	var file objects.CodeFile
	file.Lines = ReadyLines(strings.Split(fs.ReadAllText(path), "\n"))
	file.Path = path
	file.File = fs.OpenFile(path)
	return file
}

// ReadyLines Ready lines to process.
// lines Lines to ready.
func ReadyLines(lines []string) *vector.Vector {
	readyLines := vector.New()
	for index := 0; index < len(lines); index++ {
		readyLines.Vals = append(readyLines.Vals, objects.CodeLine{Line: index + 1,
			Text: strings.TrimRight(lines[index], " \t\n\r")})
	}
	return readyLines
}
