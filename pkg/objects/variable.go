package objects

// Variable Variable instance.
type Variable struct {
	// Name of variable.
	Name string
	// Line of define.
	Line int
	// Value of variable.
	Value Value
	// Is const variable.
	Const bool
	// Protection state.
	Protected bool
}
