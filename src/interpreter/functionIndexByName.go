/*
	functionIndexByName Function.
*/

package interpreter

import (
	"../objects"
)

// functionIndexByName Find index of function by name.
// name Name to find.
func (i Interpreter) functionIndexByName(name string) int {
	for index := range i.funcs.Vals {
		if i.funcs.Vals[index].(objects.Function).Name == name {
			return index
		}
	}
	return -1
}
