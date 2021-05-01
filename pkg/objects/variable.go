package objects

// Variable instance.
type Variable struct {
	Name      string
	Line      int     // Line of define.
	Value     Value
	Const     bool
	Protected bool
}
