package interpreter

import (
	"os"
	"strings"

	"github.com/fract-lang/fract/pkg/fs"
	"github.com/fract-lang/fract/pkg/objects"
)

// ReadyFile returns instance of source file by path.
func ReadyFile(path string) *objects.SourceFile {
	file, _ := os.Open(path)
	return &objects.SourceFile{
		Lines: ReadyLines(strings.Split(fs.ReadAllText(path), "\n")),
		Path:  path,
		File:  file,
	}
}
