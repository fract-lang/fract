/*
	processLoop
*/

package interpreter

import (
	"../fract"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// isWhile Find in keyword in tokens.
// tokens Tokens to check.
func findIn(tokens vector.Vector) int {
	for index := 0; index < tokens.Len(); index++ {
		current := tokens.At(index).(objects.Token)
		if current.Type == fract.TypeIn {
			return index
		}
	}
	return -1
}

// processLoop Process loop block.
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processLoop(tokens *vector.Vector, do bool) {
	index := parser.IndexBlockDeclare(tokens)
	// Block declare is not defined?
	if index == -1 {
		fract.Error(tokens.Last().(objects.Token), "Where is the block declare!?")
	}
	inIndex := findIn(tokens.Sublist(0, index))

	// While?
	if inIndex == -1 {
		conditionList := tokens.Sublist(1, tokens.Len()-inIndex-3)
		line := i.lexer.Line

		cacheList := tokens.Sublist(index+1, tokens.Len()-index-1)
		tokens = &cacheList

		/* Interpret/skip block. */
		for !i.lexer.Finished {
			// Skip this loop if tokens are empty.
			if !tokens.Any() {
				tokens = i.lexer.Next()
				continue
			}

			first := tokens.First().(objects.Token)
			if first.Type == fract.TypeBlockEnd { // Block is ended.
				if line == -1 {
					return
				}
				i.lexer.Line -= i.lexer.Line - line
				i.lexer.BlockCount++
				tokens = i.lexer.Next()
				continue
			}

			// Condition is true?
			if i.processCondition(&conditionList) == grammar.TRUE {
				if do {
					i.processTokens(tokens, do)
				}
			} else {
				line = -1
			}

			tokens = i.lexer.Next()
		}
	}
}
