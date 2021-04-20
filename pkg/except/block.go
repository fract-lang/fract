package except

import obj "github.com/fract-lang/fract/pkg/objects"

// Code block instance.
type Block struct {
	// Main block.
	Try func()
	// On panic catch.
	Catch func(obj.Exception)
}
