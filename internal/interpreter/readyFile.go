/*
	ReadyFile Function.
*/

package interpreter

import (
	"strings"

	"github.com/fract-lang/fract/pkg/fs"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// ReadyFile Create instance of code file.
// path Path of file.
func ReadyFile(path string) obj.CodeFile {
	return obj.CodeFile{
		Lines: ReadyLines(strings.Split(fs.ReadAllText(path), "\n")),
		Path:  path,
		File:  fs.OpenFile(path),
	}
}
