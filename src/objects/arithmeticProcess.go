package objects

// ArithmeticProcess Arithmetic process instance.
type ArithmeticProcess struct {
	// First value of process.
	First Token
	// Value instance of first value.
	FirstV Value
	// Second value of process.
	Second Token
	// Value instance of second value.
	SecondV Value
	// Operator of process.
	Operator Token
}
