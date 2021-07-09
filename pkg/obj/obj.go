package obj

import (
	"os"
)

// Func instance.
type Func struct {
	Name              string
	Ln                int      // Line of define.
	Tks               []Tokens // Block content of function.
	Params            []Param
	DefaultParamCount int
	Protected         bool
}

// Var instance.
type Var struct {
	Name      string
	Ln        int // Line of define.
	Val       Value
	Const     bool
	Mut       bool
	Protected bool
}

// Token instance.
type Token struct {
	F   *File
	Val string
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
