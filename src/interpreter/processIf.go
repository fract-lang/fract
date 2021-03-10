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

// processIf Process if-elif-else blocks and returns loop keyword state.
// And returns loop keyword state.
//
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processIf(tokens vector.Vector, do bool) int {
	i.blockCount++

	/* IF */
	tokenLen := len(tokens.Vals)

	// Block declare is not defined?
	if tokens.Vals[tokenLen-1].(objects.Token).Type != fract.TypeBlock {
		fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token),
			"Where is the block declare!?")
	}
	conditionList := tokens.Sublist(1, tokenLen-2)

	// Condition is empty?
	if len(conditionList.Vals) == 0 {
		first := tokens.Vals[0].(objects.Token)
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Condition is empty!")
	}

	state := i.processCondition(conditionList)
	actioned := state == grammar.TRUE
	kwstate := fract.TypeNone

	/* Interpret/skip block. */
	i.index++
	for ; i.index < len(i.tokens.Vals); i.index++ {
		tokens := i.tokens.Vals[i.index].(vector.Vector)
		do = kwstate == -1 && do

		first := tokens.Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			i.blockCount--
			return kwstate
		} else if first.Type == fract.TypeElseIf { // Else if block.
			tokenLen = len(tokens.Vals)

			// Block declare is not defined?
			if tokens.Vals[tokenLen-1].(objects.Token).Type != fract.TypeBlock {
				fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token),
					"Where is the block declare!?")
			}
			conditionList := tokens.Sublist(1, tokenLen-2)

			// Condition is empty?
			if len(conditionList.Vals) == 0 {
				first := tokens.Vals[0].(objects.Token)
				fract.ErrorCustom(first.File.Path, first.Line,
					first.Column+len(first.Value), "Condition is empty!")
			}

			state = i.processCondition(conditionList)

			/* Interpret/skip block. */
			i.index++
			for ; i.index < len(i.tokens.Vals); i.index++ {
				tokens := i.tokens.Vals[i.index].(vector.Vector)

				first := tokens.Vals[0].(objects.Token)
				if first.Type == fract.TypeBlockEnd { // Block is ended.
					i.blockCount--
					return kwstate
				} else if first.Type == fract.TypeIf { // If block.
					i.processIf(tokens, state == grammar.TRUE && !actioned && do)
					continue
				} else if first.Type == fract.TypeElseIf ||
					first.Type == fract.TypeElse { // Else if or else block.
					break
				}

				// Condition is true?
				if state == grammar.TRUE && !actioned && do {
					if kwstate = i.processTokens(tokens, true); kwstate != fract.TypeNone {
						i.skipBlock()
					}
				}
			}

			if state == grammar.TRUE {
				i.skipBlock()
				i.index--
			} else if !actioned {
				actioned = state == grammar.TRUE
			}
			continue
		} else if first.Type == fract.TypeElse { // Else block.
			// Block declare is not defined?
			if tokens.Vals[len(tokens.Vals)-1].(objects.Token).Type != fract.TypeBlock {
				fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token),
					"Where is the block declare!?")
			}

			/* Interpret/skip block. */
			i.index++
			for ; i.index < len(i.tokens.Vals); i.index++ {
				tokens := i.tokens.Vals[i.index].(vector.Vector)

				first := tokens.Vals[0].(objects.Token)
				if first.Type == fract.TypeBlockEnd { // Block is ended.
					i.blockCount--
					return kwstate
				} else if first.Type == fract.TypeIf { // If block.
					i.processIf(tokens, !actioned && do)
					continue
				}

				// Condition is true?
				if !actioned && do {
					if kwstate = i.processTokens(tokens, true); kwstate != fract.TypeNone {
						i.skipBlock()
					}
				}
			}
		}

		// Condition is true?
		if state == grammar.TRUE && do {
			if kwstate = i.processTokens(tokens, do); kwstate != fract.TypeNone {
				i.skipBlock()
			}
		}
	}
	return kwstate
}
