/*
	Interpret Function
*/

package interpreter

import (
	"../fract"
	"../objects"
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

	// Change blocks.
	count := 0
	for i.index = range i.tokens.Vals {
		first := i.tokens.Vals[i.index].(vector.Vector).Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd {
			count--
		} else if first.Type == fract.TypeIf || first.Type == fract.TypeLoop ||
			first.Type == fract.TypeFunction {
			count++
		}
	}

	if count > 0 { // Check blocks.
		i.lexer.Line--
		i.lexer.Error("Block is expected ending...")
	}

	// Interpret all lines.
	for i.index = 0; i.index < len(i.tokens.Vals); i.index++ {
		i.processTokens(i.tokens.Vals[i.index].(vector.Vector), true)
	}
}
