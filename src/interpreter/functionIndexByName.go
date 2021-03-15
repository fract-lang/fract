/*
	functionIndexByName Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/objects"
)

// functionIndexByName Find index of function by name.
// name Name to find.
func (i Interpreter) functionIndexByName(name string) int {
	for index, current := range i.funcs.Vals {
		if current.(objects.Function).Name == name {
			return index
		}
	}
	return -1
}
