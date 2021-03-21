/*
	RemoveRange Function.
*/

package vector

import (
	obj "github.com/fract-lang/fract/pkg/objects"
)

// RemoveRange Remove range.
// slice Source slice.
// pos Start position of removing.
// len Count of removing elements.
func RemoveRange(slice *[]obj.Token, pos, len int) {
	*slice = append((*slice)[:pos], (*slice)[pos+len:]...)
}
