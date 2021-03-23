/*
	PrintValue Function.
*/

package fract

import (
	"fmt"

	obj "github.com/fract-lang/fract/pkg/objects"
)

// PrintValue Print value to screen.
// value Value to print.
func PrintValue(value obj.Value) bool {
	if value.Content == nil {
		return false
	}

	if value.Array {
		fmt.Print(value.Content)
	} else {
		fmt.Print(value.Content[0])
	}
	return true
}
