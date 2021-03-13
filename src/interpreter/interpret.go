/*
	Interpret Function
*/

package interpreter

import (
	"../utils/vector"
)

// Interpret Interpret file.
func (i *Interpreter) Interpret() {
	// Lexer is finished.
	if i.lexer.Finished {
		return
	}

	/* Tokenize all lines. */
	for !i.lexer.Finished {
		cacheTokens := i.lexer.Next()

		// cacheTokens are empty?
		if cacheTokens.Vals == nil {
			continue
		}

		i.tokens.Vals = append(i.tokens.Vals, cacheTokens)
	}

	// Interpret all lines.
	for ; i.index < len(i.tokens.Vals); i.index++ {
		i.processTokens(i.tokens.Vals[i.index].(vector.Vector), true)
	}

	if i.blockCount > 0 { // Check blocks.
		i.lexer.Line--
		i.lexer.Error("Block is expected ending...")
	}
}
