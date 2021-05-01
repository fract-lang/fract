package interpreter

import (
	"strings"
	"unicode"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

// functionIndexByName returns index of function by name.
func (i *Interpreter) functionIndexByName(name objects.Token) (int, *Interpreter) {
	if name.Value[0] == '-' { // Ignore.
		name.Value = name.Value[1:]
	}

	if index := strings.Index(name.Value, grammar.TokenDot); index != -1 {
		if i.importIndexByName(name.Value[:index]) == -1 {
			fract.Error(name, "'"+name.Value[:index]+"' is not defined!")
		}
		i = i.Imports[i.importIndexByName(name.Value[:index])].Source
		name.Value = name.Value[index+1:]

		for index, current := range i.functions {
			if !unicode.IsUpper(rune(current.Name[0])) {
				continue
			}

			if current.Name == name.Value {
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
