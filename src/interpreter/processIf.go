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
// do Do processes?
func (i *Interpreter) processIf(tokens *vector.Vector, do bool) {
	/* IF */
	last := tokens.Last().(objects.Token)
	// Block declare is not defined?
	if last.Type != fract.TypeBlock {
		fract.Error(last, "Where is the block declare!?")
	}
	conditionList := tokens.Sublist(1, tokens.Len()-2)
	state := i.processCondition(&conditionList)
	actioned := state == grammar.TRUE

	// Go next line.
	tokens = i.lexer.Next()

	/* Interpret/skip block. */
	for !i.lexer.Finished {
		// Skip this loop if tokens are empty.
		if !tokens.Any() {
			return
		}

		first := tokens.First().(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			return
		} else if first.Type == fract.TypeIf { // If block.
			i.processIf(tokens, state == grammar.TRUE && do)
		} else if first.Type == fract.TypeElseIf { // Else if block.
			i.lexer.BlockCount--

			last = tokens.Last().(objects.Token)
			// Block declare is not defined?
			if last.Type != fract.TypeBlock {
				fract.Error(last, "Where is the block declare!?")
			}
			conditionList := tokens.Sublist(1, tokens.Len()-2)
			state = i.processCondition(&conditionList)

			tokens = i.lexer.Next()
			/* Interpret/skip block. */
			for !i.lexer.Finished {
				// Skip this loop if tokens are empty.
				if !tokens.Any() {
					return
				}

				first := tokens.First().(objects.Token)
				if first.Type == fract.TypeBlockEnd { // Block is ended.
					return
				} else if first.Type == fract.TypeIf { // If block.
					i.processIf(tokens, state == grammar.TRUE && !actioned && do)
				} else if first.Type == fract.TypeElseIf { // Else if block.
					break
				}

				// Condition is true?
				if state == grammar.TRUE && !actioned && do {
					i.processTokens(tokens)
				}

				tokens = i.lexer.Next()
			}

			if !actioned {
				actioned = state == grammar.TRUE
			}
			continue
		}

		// Condition is true?
		if state == grammar.TRUE && do {
			i.processTokens(tokens)
		}

		tokens = i.lexer.Next()
	}
}
