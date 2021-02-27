/*
	VarIndexByName Function.
*/

package name

import (
	"../../objects"
	"../../utilities/vector"
)

// VarIndexByName Find index of variable by name.
// vars All variables.
// name Name to find.
func VarIndexByName(vars *vector.Vector, name string) int {
	for index := range vars.Vals {
		if vars.Vals[index].(objects.Variable).Name == name {
			return index
		}
	}
	return -1
}
