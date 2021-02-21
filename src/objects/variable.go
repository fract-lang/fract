package objects

// Variable Variable instance.
type Variable struct {
	// Name of variable.
	Name string
	// Value of variable.
	Value []string
	// Type of variable.
	Type string
	// Is const variable.
	Const bool
	// Is array variable.
	Array bool
	// Is char typed variable.
	Charray bool
}
