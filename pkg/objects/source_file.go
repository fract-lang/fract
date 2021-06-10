package objects

import "os"

// Source file instance.
type SourceFile struct {
	Path  string
	File  *os.File
	Lines []string
}
