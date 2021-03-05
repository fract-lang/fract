package objects

import (
	"../utilities/vector"
)

// Function Function instance.
type Function struct {
	// Name of function.
	Name string
	// Block start of function.
	Start int
	// DataType of return value.
	ReturnType string
	// Parameters of function.
	Parameters *vector.Vector
}
