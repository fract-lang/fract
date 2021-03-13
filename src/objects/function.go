package objects

// Function Function instance.
type Function struct {
	// Name of function.
	Name string
	// Block start of function.
	Start int
	// Block content of function.
	Tokens []interface{}
	// Parameters of function.
	Parameters []string
}
