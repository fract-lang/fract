/*
	VarIndexByName Function.
*/

package name

import (
	"../../objects"
	"../../utilities/vector"
)

// VarIndexByName Find index of variable by name.
// name Name to find.
func VarIndexByName(vars *vector.Vector, name string) int {
	for index := range vars.Vals {
		current := vars.Vals[index].(objects.Variable)
		if current.Name == name {
			return index
		}
	}
	return -1
}
