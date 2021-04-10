/*
	Import Function.
*/

package interpreter

import (
	"unicode"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
)

// Import content into destination interpeter.
// dest Destination interpreter.
func (i *Interpreter) Import(dest *Interpreter, name string) {
	i.ready()

	// Interpret all lines.
	for i.index = 0; i.index < len(i.Tokens); i.index++ {
		tokens := i.Tokens[i.index]
		switch first := tokens[0]; first.Type {
		case fract.TypeProtected: // Protected declaration.
			if len(tokens) < 2 {
				fract.Error(first, "Protected but what is it protected?")
			}
			second := tokens[1]
			tokens = tokens[1:]
			if second.Type == fract.TypeVariable { // Variable definition.
				i.processVariableDefinition(tokens, true)
			} else if second.Type == fract.TypeFunction { // Function definition.
				i.processFunction(tokens, true)
			} else {
				fract.Error(second, "Syntax error, you can protect only deletable objects!")
			}
		case fract.TypeVariable: // Variable definition.
			i.processVariableDefinition(tokens, false)
		case fract.TypeFunction: // Function definiton.
			i.processFunction(tokens, false)
		case fract.TypeImport: // Import.
			real := dest.Dir
			dest.Dir = i.Dir
			dest.processImport(tokens)
			dest.Dir = real
		}
	}

	// Variables.
	for _, variable := range i.variables {
		if !unicode.IsUpper(rune(variable.Name[0])) {
			continue
		}

		variable.Name = name + grammar.TokenDot + variable.Name
		dest.variables = append(dest.variables, variable)
	}

	// Functions.
	for _, function := range i.functions {
		if !unicode.IsUpper(rune(function.Name[0])) {
			continue
		}

		function.Name = name + grammar.TokenDot + function.Name
		dest.functions = append(dest.functions, function)
	}
}
