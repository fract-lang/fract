package interpreter

// ImportInfo Information of import.
type ImportInfo struct {
	// Package name.
	Name string
	// Source of package.
	Source *Interpreter
}
