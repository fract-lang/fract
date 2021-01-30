package interpreter

import (
	"../objects"
	"../utilities/fs"
)

// Interprater Interprater of Fract.
type Interprater struct {
	/* PRIVATE */

	// Parser of this file.
	file objects.CodeFile

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
func New(path string, _type int) *Interprater {
	var parser *Interprater = new(Interprater)
	parser.file = ReadyFile(path)
	parser.Type = _type
	return parser
}
