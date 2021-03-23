/*
	processIf Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processIf Process if-elif-else blocks and returns loop keyword state.
// And returns loop keyword state.
//
// tokens Tokens to process.
func (i *Interpreter) processIf(tokens []obj.Token) int {
	/* IF */
	tokenLen := len(tokens)
	conditionList := *vector.Sublist(tokens, 1, tokenLen-1)

	// Condition is empty?
	if conditionList == nil {
		first := tokens[0]
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value),
			"Condition is empty!")
	}

	state := i.processCondition(&conditionList)
	actioned := state == grammar.KwTrue
	variableLen := len(i.vars)
	functionLen := len(i.funcs)
	kwstate := fract.TypeNone

	/* Interpret/skip block. */
	for {
		i.index++
		tokens := i.Tokens[i.index]
		first := tokens[0]

		if first.Type == fract.TypeBlockEnd { // Block is ended.
			goto ret
		} else if first.Type == fract.TypeElseIf { // Else if block.
			tokenLen = len(tokens)
			conditionList := vector.Sublist(tokens, 1, tokenLen-1)

			// Condition is empty?
			if conditionList == nil {
				first := tokens[0]
				fract.ErrorCustom(first.File, first.Line,
					first.Column+len(first.Value), "Condition is empty!")
			}

			state = i.processCondition(conditionList)

			/* Interpret/skip block. */
			for {
				i.index++
				tokens := i.Tokens[i.index]
				first := tokens[0]

				if first.Type == fract.TypeBlockEnd { // Block is ended.
					goto ret
				} else if first.Type == fract.TypeIf { // If block.
					if state == grammar.KwTrue && !actioned && kwstate == fract.TypeNone {
						i.processIf(tokens)
					} else {
						i.skipBlock(true)
					}
					continue
				} else if first.Type == fract.TypeElseIf ||
					first.Type == fract.TypeElse { // Else if or else block.
					break
				}

				// Condition is true?
				if state == grammar.KwTrue && !actioned && kwstate == fract.TypeNone {
					if kwstate = i.processTokens(tokens); kwstate != fract.TypeNone {
						i.skipBlock(false)
					}
				} else {
					i.skipBlock(true)
				}
			}

			if state == grammar.KwTrue {
				i.skipBlock(false)
				i.index--
			} else if !actioned {
				actioned = state == grammar.KwTrue
			}
			continue
		} else if first.Type == fract.TypeElse { // Else block.
			if len(tokens) > 1 {
				fract.Error(first, "Else block is not take any arguments!")
			}

			/* Interpret/skip block. */
			for {
				i.index++
				tokens := i.Tokens[i.index]
				first := tokens[0]

				if first.Type == fract.TypeBlockEnd { // Block is ended.
					goto ret
				} else if first.Type == fract.TypeIf { // If block.
					if !actioned && kwstate == fract.TypeNone {
						i.processIf(tokens)
					} else {
						i.skipBlock(true)
					}
					continue
				}

				// Condition is true?
				if !actioned && kwstate == fract.TypeNone {
					if kwstate = i.processTokens(tokens); kwstate != fract.TypeNone {
						i.skipBlock(false)
					}
				}
			}
		}

		// Condition is true?
		if state == grammar.KwTrue && kwstate == fract.TypeNone {
			if kwstate = i.processTokens(tokens); kwstate != fract.TypeNone {
				i.skipBlock(false)
			}
		} else {
			i.skipBlock(true)
		}
	}
ret:
	i.vars = i.vars[:variableLen]
	i.funcs = i.funcs[:functionLen]
	return kwstate
}