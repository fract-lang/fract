package objects

import "os"

type Exception struct {
	Message string
}

// Function instance.
type Function struct {
	Name                  string
	Line                  int       // Line of define.
	Tokens                [][]Token // Block content of function.
	Parameters            []Parameter
	DefaultParameterCount int
	Protected             bool
}

// Variable instance.
type Variable struct {
	Name      string
	Line      int // Line of define.
	Value     Value
	Const     bool
	Mutable   bool
	Protected bool
}

// Token instance.
type Token struct {
	File   *SourceFile
	Value  string
	Type   uint8
	Line   int
	Column int
}

// Source file instance.
type SourceFile struct {
	Path  string
	File  *os.File
	Lines []string
}

// Parameter instance.
type Parameter struct {
	Default Value
	Name    string
	Params  bool
}
