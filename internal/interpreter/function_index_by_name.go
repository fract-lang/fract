package interpreter

import (
	"strings"
	"unicode"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

//! This code block very like to varIndexByName function. If you change here, probably you must change there too.

// functionIndexByName returns index of function by name.
func (i *Interpreter) functionIndexByName(name objects.Token) (int, *Interpreter) {
	if name.Value[0] == '-' { // Ignore minus.
		name.Value = name.Value[1:]
	}
	if index := strings.Index(name.Value, "."); index != -1 {
		if i.importIndexByName(name.Value[:index]) == -1 {
			fract.Error(name, "'"+name.Value[:index]+"' is not defined!")
		}
		i = i.Imports[i.importIndexByName(name.Value[:index])].Source
		name.Value = name.Value[index+1:]
		for index, current := range i.functions {
			if (current.Tokens == nil || unicode.IsUpper(rune(current.Name[0]))) && current.Name == name.Value {
				return index, i
			}
		}
		return -1, nil
	}
	for index, current := range i.functions {
		if current.Name == name.Value {
			return index, i
		}
	}
	return -1, nil
}
