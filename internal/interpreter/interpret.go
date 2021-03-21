/*
	Interpret Function
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
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
	{
		count := 0
		last := -1
		for i.index = range i.tokens.Vals {
			tokens := i.tokens.Vals[i.index].(vector.Vector)

			if first := tokens.Vals[0].(obj.Token); first.Type == fract.TypeBlockEnd {
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
			fract.Error(i.tokens.Vals[last].(vector.Vector).Vals[0].(obj.Token),
				"Block is expected ending...")
		}
	}

	// Interpret all lines.
	for i.index = 0; i.index < len(i.tokens.Vals); i.index++ {
		i.processTokens(i.tokens.Vals[i.index].(vector.Vector))
	}
}
