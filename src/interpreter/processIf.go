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

// processIf Process if-elif-else blocks and returns loop keyword state.
// And returns loop keyword state.
//
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processIf(tokens *vector.Vector, do bool) int {
	/* IF */
	index := parser.IndexBlockDeclare(tokens)
	// Block declare is not defined?
	if index == -1 {
		fract.Error(tokens.Last().(objects.Token), "Where is the block declare!?")
	}
	conditionList := tokens.Sublist(1, index-1)

	// Condition is empty?
	if !conditionList.Any() {
		first := tokens.First().(objects.Token)
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Condition is empty!")
	}

	state := i.processCondition(&conditionList)
	actioned := state == grammar.TRUE

	// Get after block tokens with used @conditionList as cache.
	conditionList = tokens.Sublist(index+1, tokens.Len()-index-1)
	tokens = &conditionList

	kwstate := -1

	/* Interpret/skip block. */
	for !i.lexer.Finished {
		do = kwstate == -1 && do

		// Skip this loop if tokens are empty.
		if !tokens.Any() {
			tokens = i.lexer.Next()
			continue
		}

		first := tokens.First().(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			return kwstate
		} else if first.Type == fract.TypeElseIf { // Else if block.
			i.lexer.BlockCount--

			index = parser.IndexBlockDeclare(tokens)
			// Block declare is not defined?
			if index == -1 {
				fract.Error(tokens.Last().(objects.Token), "Where is the block declare!?")
			}
			conditionList := tokens.Sublist(1, index-1)

			// Condition is empty?
			if !conditionList.Any() {
				first := tokens.First().(objects.Token)
				fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
					"Condition is empty!")
			}

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
					return kwstate
				} else if first.Type == fract.TypeIf { // If block.
					i.processIf(tokens, state == grammar.TRUE && !actioned && do)
				} else if first.Type == fract.TypeElseIf { // Else if block.
					break
				}

				// Condition is true?
				if state == grammar.TRUE && !actioned && do {
					kwstate = i.processTokens(tokens, true)
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
			kwstate = i.processTokens(tokens, do)
		}

		tokens = i.lexer.Next()
	}
	return kwstate
}
