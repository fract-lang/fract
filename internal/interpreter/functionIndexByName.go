/*
	functionIndexByName Function.
*/

package interpreter

// functionIndexByName Find index of function by name.
// name Name to find.
func (i Interpreter) functionIndexByName(name string) int {
	for index, current := range i.functions {
		if current.Name == name {
			return index
		}
	}
	return -1
}
