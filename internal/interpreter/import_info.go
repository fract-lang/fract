package interpreter

// ImportInfo Information of import.
type ImportInfo struct {
	Name   string         // Package name.
	Source *Interpreter   // Source of package.
}
