/*
	varIndexByName Function.
*/

package interpreter

import (
	"../objects"
)

// varIndexByName Find index of variable by name.
// name Name to find.
func (i *Interpreter) varIndexByName(name string) int {
	for index := range i.vars.Vals {
		if i.vars.Vals[index].(objects.Variable).Name == name {
			return index
		}
	}
	return -1
}
