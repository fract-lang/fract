package objects

// Token Token instance.
type Token struct {
	// File of token.
	File CodeFile
	// Value of token.
	Value string
	// Type of token.
	Type int
	// Line of token.
	Line int
	// Column of token.
	Column int
}
