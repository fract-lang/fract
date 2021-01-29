package objects

import (
	"container/list"
	"os"
)

// CodeFile Code file instance.
type CodeFile struct {
	// Path of file.
	path string
	// File instance of file.
	file *os.File
	// Lines of file.
	lines list.List
}
