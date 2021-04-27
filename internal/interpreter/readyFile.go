/*
	ReadyFile Function.
*/

package interpreter

import (
	"os"
	"strings"

	"github.com/fract-lang/fract/pkg/fs"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// ReadyFile Create instance of code file.
// path Path of file.
func ReadyFile(path string) *obj.CodeFile {
	file, _ := os.Open(path)
	return &obj.CodeFile{
		Lines: ReadyLines(strings.Split(fs.ReadAllText(path), "\n")),
		Path:  path,
		File:  file,
	}
}
