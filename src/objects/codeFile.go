package objects

import (
	"os"

	"github.com/fract-lang/fract/src/utils/vector"
)

// CodeFile Code file instance.
type CodeFile struct {
	// Path of file.
	Path string
	// File instance of file.
	File *os.File
	// Lines of file.
	Lines vector.Vector
}
