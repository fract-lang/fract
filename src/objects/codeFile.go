package objects

import (
	"os"

	"../utilities/vector"
)

// CodeFile Code file instance.
type CodeFile struct {
	// Path of file.
	Path string
	// File instance of file.
	File *os.File
	// Lines of file.
	Lines *vector.Vector
}
