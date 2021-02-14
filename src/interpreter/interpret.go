/*
	Interpret Function
*/

package interpreter

import "../utilities/vector"

// Interpret Interpret file.
func (i *Interpreter) Interpret() {
	// Lexer is finished.
	if i.lexer.Finished {
		return
	}

	/* Interpret all lines. */
	for !i.lexer.Finished {
		cacheTokens := i.lexer.Next()

		// cacheTokens are empty?
		if !cacheTokens.Any() {
			continue
		}

		i.tokens.Append(cacheTokens)
	}

	i.tokenLen = i.tokens.Len()
	for ; i.index < i.tokenLen; i.index++ {
		i.processTokens(i.tokens.At(i.index).(*vector.Vector), true)
	}

	if i.blockCount > 0 { // Check blocks.
		i.lexer.Line--
		i.lexer.Error("Block is expected ending...")
	}
}
