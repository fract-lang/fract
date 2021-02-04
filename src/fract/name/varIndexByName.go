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
	for index := 0; index < len(vars.Vals); index++ {
		current := vars.At(index).(objects.Variable)
		if current.Name == name {
			return index
		}
	}
	return -1
}
