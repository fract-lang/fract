package interpreter

// importIndexByName returns index of import by name.
func (i *Interpreter) importIndexByName(name string) int {
	for index, current := range i.Imports {
		if current.Name == name {
			return index
		}
	}
	return -1
}
