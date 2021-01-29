package objects

import (
	"container/list"
	"os"
)

// CodeFile Code file instance.
type CodeFile struct {
	// Path of file.
	Path string
	// File instance of file.
	File *os.File
	// Lines of file.
	Lines list.List
}
