/*
	functionIndexByName Function.
*/

package interpreter

import (
	"strings"
	"unicode"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
)

// functionIndexByName Find index of function by name.
// name Name to find.
func (i *Interpreter) functionIndexByName(name obj.Token) (int, *Interpreter) {
	index := strings.Index(name.Value, grammar.TokenDot)

	if index != -1 {
		iindex := i.importIndexByName(name.Value[:index])
		if iindex == -1 {
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
