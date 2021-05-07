package interpreter

import (
	"os"
	"strings"

	"github.com/fract-lang/fract/pkg/fs"
	"github.com/fract-lang/fract/pkg/objects"
)

// ReadyFile returns instance of source file by path.
func ReadyFile(filepath string) *objects.SourceFile {
	file, _ := os.Open(filepath)
	return &objects.SourceFile{
		Lines: ReadyLines(strings.Split(fs.ReadAllText(filepath), "\n")),
		Path:  filepath,
		File:  file,
	}
}
