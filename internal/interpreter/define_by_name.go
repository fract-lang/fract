package interpreter

import "github.com/fract-lang/fract/pkg/objects"

// TYPES
// 'f' -> Function.
// 'v' -> Variable.

func (i *Interpreter) defineByName(name objects.Token) (int, rune, *Interpreter) {
	index, source := i.functionIndexByName(name)
	if index != -1 {
		return index, 'f', source
	}
	index, source = i.variableIndexByName(name)
	if index != -1 {
		return index, 'v', source
	}
	return -1, '-', nil
}

func (i *Interpreter) definedName(name objects.Token) int {
	if name.Value[0] == '-' { // Ignore minus.
		name.Value = name.Value[1:]
	}
	for _, current := range i.functions {
		if current.Name == name.Value {
			return current.Line
		}
	}
	for _, current := range i.variables {
		if current.Name == name.Value {
			return current.Line
		}
	}
	return -1
}
