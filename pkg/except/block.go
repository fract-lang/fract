package except

import "github.com/fract-lang/fract/pkg/objects"

// Code block instance.
type Block struct {
	Try       func()
	Catch     func(*objects.Exception)
	Exception *objects.Exception
}
