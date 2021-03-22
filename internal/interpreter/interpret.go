/*
	Interpret Function
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/parser"
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
		if cacheTokens == nil {
			continue
		}

		i.tokens = append(i.tokens, cacheTokens)
	}

	// Change blocks.
	count := 0
	last := -1
	for i.index = range i.tokens {
		tokens := i.tokens[i.index]

		if first := tokens[0]; first.Type == fract.TypeBlockEnd {
			count--
			if count < 0 {
				fract.Error(first, "The extra block end defined!")
			}
		} else if parser.IsBlockStatement(tokens) {
			count++
			if count == 1 {
				last = i.index
			}
		}
	}

	if count > 0 { // Check blocks.
		fract.Error(i.tokens[last][0],
			"Block is expected ending...")
	}

	// Interpret all lines.
	for i.index = 0; i.index < len(i.tokens); i.index++ {
		i.processTokens(i.tokens[i.index])
	}
}
