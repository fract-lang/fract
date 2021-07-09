package obj

import (
	"os"
)

// Var instance.
type Var struct {
	Name      string
	Ln        int // Line of define.
	V         Value
	Const     bool
	Mut       bool
	Protected bool
}

// Token instance.
type Token struct {
	F   *File
	V   string
	T   uint8
	Ln  int
	Col int
}

// Source file instance.
type File struct {
	P   string
	F   *os.File
	Lns []string
}

// Param instance.
type Param struct {
	Default Value
	Name    string
	Params  bool
}
