/*
	varIndexByName Function.
*/

package interpreter

// varIndexByName Find index of variable by name.
// name Name to find.
func (i Interpreter) varIndexByName(name string) int {
	for index, current := range i.variables {
		if current.Name == name {
			return index
		}
	}
	return -1
}
