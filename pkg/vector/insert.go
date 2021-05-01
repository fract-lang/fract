package vector

import "github.com/fract-lang/fract/pkg/objects"

// Insert value by position.
func Insert(slice *[]objects.Token, pos int, value ...objects.Token) {
	*slice = append((*slice)[:pos], append(value, (*slice)[pos:]...)...)
}
