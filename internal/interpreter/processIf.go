package interpreter

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processIf process if-elif-else blocks and returns keyword state.
func (i *Interpreter) processIf(tokens []objects.Token) uint8 {
	tokenLen := len(tokens)
	conditionList := vector.Sublist(tokens, 1, tokenLen-1)

	// Condition is empty?
	if conditionList == nil {
		first := tokens[0]
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value), "Condition is empty!")
	}

	state := i.processCondition(conditionList)
	variableLen := len(i.variables)
	functionLen := len(i.functions)
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

			if state == grammar.KwTrue {
				i.skipBlock(false)
				goto ret
			}

			state = i.processCondition(conditionList)

			// Interpret/skip block.
			for {
				i.index++
				tokens := i.Tokens[i.index]
				first := tokens[0]

				if first.Type == fract.TypeBlockEnd { // Block is ended.
					goto ret
				} else if first.Type == fract.TypeIf { // If block.
					if state == grammar.KwTrue && kwstate == fract.TypeNone {
						i.processIf(tokens)
					} else {
						i.skipBlock(true)
					}
					continue
				} else if first.Type == fract.TypeElseIf || first.Type == fract.TypeElse { // Else if or else block.
					i.index--
					break
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

			if state == grammar.KwTrue {
				i.skipBlock(false)
				goto ret
			}
			continue
		} else if first.Type == fract.TypeElse { // Else block.
			if len(tokens) > 1 {
				fract.Error(first, "Else block is not take any arguments!")
			}

			if state == grammar.KwTrue {
				i.skipBlock(false)
				goto ret
			}

			/* Interpret/skip block. */
			for {
				i.index++
				tokens := i.Tokens[i.index]
				first := tokens[0]

				if first.Type == fract.TypeBlockEnd { // Block is ended.
					goto ret
				} else if first.Type == fract.TypeIf { // If block.
					if kwstate == fract.TypeNone {
						i.processIf(tokens)
					} else {
						i.skipBlock(true)
					}
					continue
				}

				// Condition is true?
				if kwstate == fract.TypeNone {
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
	i.variables = i.variables[:variableLen]
	i.functions = i.functions[:functionLen]
	return kwstate
}
