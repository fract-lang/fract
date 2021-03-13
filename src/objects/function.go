package objects

import (
	"../utils/vector"
)

// Function Function instance.
type Function struct {
	// Name of function.
	Name string
	// Block start of function.
	Start int
	// Block content of function.
	Tokens vector.Vector
	// Parameters of function.
	Parameters []string
}
