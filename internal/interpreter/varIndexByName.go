/*
	varIndexByName Function.
*/

package interpreter

import (
	"strings"
	"unicode"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// varIndexByName Find index of variable by name.
// name Name to find.
func (i *Interpreter) varIndexByName(name obj.Token) (int, *Interpreter) {
	if name.Value[0] == '-' { // Ignore.
		name.Value = name.Value[1:]
	}

	index := strings.Index(name.Value, grammar.TokenDot)

	if index != -1 {
		iindex := i.importIndexByName(name.Value[:index])
		if iindex == -1 {
			fract.Error(name, "'"+name.Value[:index]+"' is not defined!")
		}
		i = i.Imports[iindex].Source
		name.Value = name.Value[index+1:]

		for index, current := range i.variables {
			if !unicode.IsUpper(rune(current.Name[0])) {
				continue
			}

			if current.Name == name.Value {
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
