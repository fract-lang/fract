/*
	Insert Function.
*/

package vector

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Insert Insert value by position.
// slice Source slice.
// pos Position to insert.
// value Value to insert.
func Insert(slice *[]obj.Token, pos int, value ...obj.Token) {
	*slice = append((*slice)[:pos], append(value, (*slice)[pos:]...)...)
}
