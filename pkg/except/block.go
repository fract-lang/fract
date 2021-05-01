package except

import obj "github.com/fract-lang/fract/pkg/objects"

// Code block instance.
type Block struct {
	Try       func()
	Catch     func(*obj.Exception)
	Exception *obj.Exception
}
