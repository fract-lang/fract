package vector

import "github.com/fract-lang/fract/pkg/objects"

// Sublist returns slice by position and length.
func Sublist(slice []objects.Token, position, length int) *[]objects.Token {
	if length == 0 {
		return nil
	}
	slice = append([]objects.Token{}, slice[position:position+length]...)
	return &slice
}

// RemoveRange by position and length.
func RemoveRange(slice *[]objects.Token, position, length int) {
	if length > 0 {
		*slice = append((*slice)[:position], (*slice)[position+length:]...)
	}
}

// Insert value by position.
func Insert(slice *[]objects.Token, pos int, value ...objects.Token) {
	*slice = append((*slice)[:pos], append(value, (*slice)[pos:]...)...)
}
