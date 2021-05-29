package interpreter

import "github.com/fract-lang/fract/pkg/fract"

// Import content into destination interpeter.
func (i *Interpreter) Import() {
	i.ready()
	// Interpret all lines.
	for i.index = 0; i.index < len(i.Tokens); i.index++ {
		switch tokens := i.Tokens[i.index]; tokens[0].Type {
		case fract.TypeProtected: // Protected declaration.
			if len(tokens) < 2 {
				fract.Error(tokens[0], "Protected but what is it protected?")
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
			source := new(Interpreter)
			source.ApplyEmbedFunctions()
			source.processImport(tokens)

			i.variables = append(i.variables, source.variables...)
			i.functions = append(i.functions, source.functions[:]...)
			i.Imports = append(i.Imports, source.Imports...)
		case fract.TypeMacro: // Macro.
			i.processMacro(tokens)
		default:
			i.skipBlock(true)
		}
	}
}
