package objects

// Function Function instance.
type Function struct {
	// Name of function.
	Name string
	// Line of define.
	Line int
	// Block content of function.
	Tokens [][]Token
	// Parameters of function.
	Parameters []Parameter
	// Count of parameters with default value.
	DefaultParameterCount int
	// Protection state?
	Protected bool
}
