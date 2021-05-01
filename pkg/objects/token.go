package objects

// Token instance.
type Token struct {
	File   *SourceFile
	Value  string
	Type   uint8
	Line   int
	Column int
}
