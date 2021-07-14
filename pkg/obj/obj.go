package obj

import (
	"os"

	"github.com/fract-lang/fract/pkg/value"
)

// Var instance.
type Var struct {
	Name      string
	Ln        int // Line of define.
	V         value.Val
	Const     bool
	Mut       bool
	Protected bool
}

// Source file instance.
type File struct {
	P   string
	F   *os.File
	Lns []string
}
