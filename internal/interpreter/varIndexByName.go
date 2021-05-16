package interpreter

import (
	"strings"
	"unicode"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
)

//! This code block very like to functionIndexByName function. If you change here, probably you must change there too.

// varIndexByName returns index of variable by name.
func (i *Interpreter) varIndexByName(name objects.Token) (int, *Interpreter) {
	if name.Value[0] == '-' { // Ignore minus.
		name.Value = name.Value[1:]
	}

	if index := strings.Index(name.Value, grammar.TokenDot); index != -1 {
		if iindex := i.importIndexByName(name.Value[:index]); iindex == -1 {
			fract.Error(name, "'"+name.Value[:index]+"' is not defined!")
		} else {
			i = i.Imports[iindex].Source
		}
		name.Value = name.Value[index+1:]

		for index, current := range i.variables {
			if (current.Line == -1 || unicode.IsUpper(rune(current.Name[0]))) && current.Name == name.Value {
				return index, i
			}
		}
		return -1, nil
	}

	for index, current := range i.variables {
		if current.Name == name.Value {
			return index, i
		}
	}
	return -1, nil
}
