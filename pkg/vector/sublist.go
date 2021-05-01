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
