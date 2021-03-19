package objects

// Function Function instance.
type Function struct {
	// Name of function.
	Name string
	// Block start of function.
	Start int
	// Line of define.
	Line int
	// Block content of function.
	Tokens []interface{}
	// Parameters of function.
	Parameters []Parameter
	// Count of parameters with default value.
	DefaultParameterCount int
}
