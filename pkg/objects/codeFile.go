package objects

import (
	"os"
)

// Code file instance.
type CodeFile struct {
	Path  string
	File  *os.File
	Lines []string
}
