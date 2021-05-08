package objects

// Function instance.
type Function struct {
	Name                  string
	Line                  int       // Line of define.
	Tokens                [][]Token // Block content of function.
	Parameters            []Parameter
	DefaultParameterCount int
	Protected             bool
}
