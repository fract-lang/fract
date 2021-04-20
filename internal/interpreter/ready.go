package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/parser"
)

// ready Ready to process.
func (i *Interpreter) ready() {
	/* Tokenize all lines. */
	for !i.Lexer.Finished {
		cacheTokens := i.Lexer.Next()

		// cacheTokens are empty?
		if cacheTokens == nil {
			continue
		}

		i.Tokens = append(i.Tokens, cacheTokens)
	}

	// Change blocks.
	count := 0
	last := -1
	for i.index = range i.Tokens {
		tokens := i.Tokens[i.index]

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
		fract.Error(i.Tokens[last][0], "Block is expected ending...")
	}
}
