/*
	Interpret Function
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
)

// Interpret Interpret file.
func (i *Interpreter) Interpret() {
	if i.Lexer.File.Path == fract.Stdin {
		// Interpret all lines.
		for i.index = 0; i.index < len(i.Tokens); i.index++ {
			i.processTokens(i.Tokens[i.index])
		}
		return
	}

	// Lexer is finished.
	if i.Lexer.Finished {
		return
	}

	i.ready()

	// Interpret all lines.
	for i.index = 0; i.index < len(i.Tokens); i.index++ {
		i.processTokens(i.Tokens[i.index])
	}
}
