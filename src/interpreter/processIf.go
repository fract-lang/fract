/*
	processIf Function.
*/

package interpreter

import (
	"../fract"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processIf Process if-elif-else blocks.
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processIf(tokens *vector.Vector, do bool) {
	/* IF */
	index := parser.IndexBlockDeclare(tokens)
	// Block declare is not defined?
	if index == -1 {
		fract.Error(tokens.Last().(objects.Token), "Where is the block declare!?")
	}
	conditionList := tokens.Sublist(1, index-1)
	state := i.processCondition(&conditionList)
	actioned := state == grammar.TRUE

	// Get after block tokens with used @conditionList as cache.
	conditionList = tokens.Sublist(index+1, tokens.Len()-index-1)
	tokens = &conditionList

	/* Interpret/skip block. */
	for !i.lexer.Finished {
		// Skip this loop if tokens are empty.
		if !tokens.Any() {
			tokens = i.lexer.Next()
			continue
		}

		first := tokens.First().(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			return
		} else if first.Type == fract.TypeIf { // If block.
			i.processIf(tokens, state == grammar.TRUE && do)
		} else if first.Type == fract.TypeElseIf { // Else if block.
			i.lexer.BlockCount--

			index = parser.IndexBlockDeclare(tokens)
			// Block declare is not defined?
			if index == -1 {
				fract.Error(tokens.Last().(objects.Token), "Where is the block declare!?")
			}
			conditionList := tokens.Sublist(1, index-1)
			state = i.processCondition(&conditionList)

			// Get after block tokens with used @conditionList as cache.
			conditionList = tokens.Sublist(index+1, tokens.Len()-index-1)
			tokens = &conditionList

			/* Interpret/skip block. */
			for !i.lexer.Finished {
				// Skip this loop if tokens are empty.
				if !tokens.Any() {
					tokens = i.lexer.Next()
					continue
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
