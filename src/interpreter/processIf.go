/*
	processIf Function.
*/

package interpreter

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// processIf Process if-elif-else blocks.
// tokens Tokens to process.
func (i *Interpreter) processIf(tokens *vector.Vector) {
	/* IF*/
	last := tokens.Last().(objects.Token)
	// Block declare is not defined?
	if last.Type != fract.TypeBlock {
		fract.Error(last, "Where is the block declare!?")
	}
	conditionList := tokens.Sublist(1, tokens.Len()-2)
	state := i.processCondition(&conditionList)

	// Go next line.
	tokens = i.lexer.Next()

	/* Interpret/skip block. */
	for !i.lexer.Finished {
		// Skip this loop if tokens are empty.
		if !tokens.Any() {
			return
		}

		// Check block is ended?
		first := tokens.First().(objects.Token)
		if first.Type == fract.TypeBlockEnd {
			return
		}

		// Condition is true?
		if state == grammar.TRUE {
			i.processTokens(tokens)
		}

		tokens = i.lexer.Next()
	}
}
