package objects

// Token Token instance.
type Token struct {
	// Value of token.
	value string
	// Type of token.
	group int
	// Line of token.
	line int
	// Column of token.
	column int
}
