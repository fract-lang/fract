/*
	Interpret Function
*/

package interpreter

// Interpret Interpret file.
func (i *Interpreter) Interpret() {
	// Lexer is finished.
	if i.lexer.Finished {
		return
	}

	/* Interpret all lines. */
	for !i.lexer.Finished {
		tokens := i.lexer.Next()
		i.processTokens(tokens)
	}
}
