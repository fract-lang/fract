package except

// Code block instance.
type Block struct {
	// Main block.
	Try func()
	// On panic catch.
	Catch func(Exception)
}
