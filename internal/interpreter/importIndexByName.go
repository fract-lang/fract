/*
	importIndexByName Function.
*/

package interpreter

// importIndexByName Find index of import by name.
// name Name to find.
func (i Interpreter) importIndexByName(name string) int {
	for index, current := range i.Imports {
		if current.Name == name {
			return index
		}
	}
	return -1
}
