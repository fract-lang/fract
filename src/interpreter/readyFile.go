/*
	ReadyFile Function.
*/

package interpreter

import (
	"strings"

	"../objects"
	"../utils/fs"
)

// ReadyFile Create instance of code file.
// path Path of file.
func ReadyFile(path string) objects.CodeFile {
	var file objects.CodeFile
	file.Lines = ReadyLines(strings.Split(fs.ReadAllText(path), "\n"))
	file.Path = path
	file.File = fs.OpenFile(path)
	return file
}
