/*
	FunctionIndexByName Function.
*/

package name

import (
	"../../objects"
	"../../utilities/vector"
)

// FunctionIndexByName Find index of function by name.
// funcs All functions.
// name Name to find.
func FunctionIndexByName(funcs vector.Vector, name string) int {
	for index := range funcs.Vals {
		if funcs.Vals[index].(objects.Function).Name == name {
			return index
		}
	}
	return -1
}
