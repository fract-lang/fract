package vector

import "github.com/fract-lang/fract/pkg/objects"

// RemoveRange by position and length.
func RemoveRange(slice *[]objects.Token, position, length int) {
	if length > 0 {
		*slice = append((*slice)[:position], (*slice)[position+length:]...)
	}
}
