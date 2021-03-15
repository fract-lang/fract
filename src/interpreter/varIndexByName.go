/*
	varIndexByName Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/objects"
)

// varIndexByName Find index of variable by name.
// name Name to find.
func (i Interpreter) varIndexByName(name string) int {
	for index, current := range i.vars.Vals {
		if current.(objects.Variable).Name == name {
			return index
		}
	}
	return -1
}
