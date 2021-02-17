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
	i.blockCount++

	/* IF */
	index := parser.IndexBlockDeclare(tokens)
	// Block declare is not defined?
	if index == -1 {
		fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token), "Where is the block declare!?")
	}
	conditionList := tokens.Sublist(1, index-1)

	// Condition is empty?
	if len(conditionList.Vals) == 0 {
		first := tokens.Vals[0].(objects.Token)
		fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
			"Condition is empty!")
	}

	state := i.processCondition(conditionList)
	actioned := state == grammar.TRUE

	tokens = tokens.Sublist(index+1, len(tokens.Vals)-index-1)

	i.emptyControl(&tokens)
	kwstate := -1

	/* Interpret/skip block. */
	i.index++
	for ; i.index < len(i.tokens.Vals); i.index++ {
		tokens = i.tokens.Vals[i.index].(*vector.Vector)
		do = kwstate == -1 && do

		first := tokens.Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			i.subtractBlock(&first)
			return kwstate
		} else if first.Type == fract.TypeElseIf { // Else if block.

			index = parser.IndexBlockDeclare(tokens)
			// Block declare is not defined?
			if index == -1 {
				fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token), "Where is the block declare!?")
			}
			conditionList := tokens.Sublist(1, index-1)

			// Condition is empty?
			if len(conditionList.Vals) == 0 {
				first := tokens.Vals[0].(objects.Token)
				fract.ErrorCustom(first.File.Path, first.Line, first.Column+len(first.Value),
					"Condition is empty!")
			}

			state = i.processCondition(conditionList)
			tokens = tokens.Sublist(index+1, len(tokens.Vals)-index-1)
			i.emptyControl(&tokens)

			/* Interpret/skip block. */
			i.index++
			for ; i.index < len(i.tokens.Vals); i.index++ {
				tokens = i.tokens.Vals[i.index].(*vector.Vector)

				first := tokens.Vals[0].(objects.Token)
				if first.Type == fract.TypeBlockEnd { // Block is ended.
					i.subtractBlock(&first)
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
	}
	return kwstate
}
