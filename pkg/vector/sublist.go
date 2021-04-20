/*
	Sublist Function.
*/

package vector

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// Sublist Get range.
// slice Source slice.
// pos Start position to take.
// length Count of taken elements.
func Sublist(slice []obj.Token, pos, length int) *[]obj.Token {
	if length == 0 {
		return nil
	}
	slice = append(make([]obj.Token, 0), slice[pos:pos+length]...)
	return &slice
}
