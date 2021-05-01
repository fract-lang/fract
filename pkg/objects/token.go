package objects

// Token instance.
type Token struct {
	File   *CodeFile
	Value  string
	Type   uint8
	Line   int
	Column int
}
