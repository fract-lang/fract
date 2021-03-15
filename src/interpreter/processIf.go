/*
	processIf Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/grammar"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// processIf Process if-elif-else blocks and returns loop keyword state.
// And returns loop keyword state.
//
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processIf(tokens vector.Vector, do bool) int {
	/* IF */
	tokenLen := len(tokens.Vals)
	conditionList := tokens.Sublist(1, tokenLen-1)

	// Condition is empty?
	if conditionList.Vals == nil {
		first := tokens.Vals[0].(objects.Token)
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value),
			"Condition is empty!")
	}

	state := i.processCondition(conditionList)
	actioned := state == grammar.TRUE
	variableLen := len(i.vars.Vals)
	functionLen := len(i.funcs.Vals)
	kwstate := fract.TypeNone

	/* Interpret/skip block. */
	for {
		i.index++
		tokens := i.tokens.Vals[i.index].(vector.Vector)
		first := tokens.Vals[0].(objects.Token)
		do = kwstate == -1 && do

		if first.Type == fract.TypeBlockEnd { // Block is ended.
			goto ret
		} else if first.Type == fract.TypeElseIf { // Else if block.
			tokenLen = len(tokens.Vals)
			conditionList := tokens.Sublist(1, tokenLen-1)

			// Condition is empty?
			if conditionList.Vals == nil {
				first := tokens.Vals[0].(objects.Token)
				fract.ErrorCustom(first.File, first.Line,
					first.Column+len(first.Value), "Condition is empty!")
			}

			state = i.processCondition(conditionList)

			/* Interpret/skip block. */
			for {
				i.index++
				tokens := i.tokens.Vals[i.index].(vector.Vector)
				first := tokens.Vals[0].(objects.Token)

				if first.Type == fract.TypeBlockEnd { // Block is ended.
					goto ret
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
			if len(tokens.Vals) > 1 {
				fract.Error(first, "Else block is not take any arguments!")
			}

			/* Interpret/skip block. */
			for {
				i.index++
				tokens := i.tokens.Vals[i.index].(vector.Vector)
				first := tokens.Vals[0].(objects.Token)

				if first.Type == fract.TypeBlockEnd { // Block is ended.
					goto ret
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
ret:
	i.vars.Vals = i.vars.Vals[:variableLen]
	i.funcs.Vals = i.funcs.Vals[:functionLen]
	return kwstate
}
