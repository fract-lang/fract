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

// processLoop Process loop block.
// tokens Tokens to process.
// do Do processes?
func (i *Interpreter) processLoop(tokens *vector.Vector, do bool) {
	index := parser.IndexBlockDeclare(tokens)
	// Block declare is not defined?
	if index == -1 {
		fract.Error(tokens.Last().(objects.Token), "Where is the block declare!?")
	}

	contentList := tokens.Sublist(1, index-1)
	// Content is empty?
	if !contentList.Any() {
		fract.Error(tokens.First().(objects.Token), "Content is empty!")
	}

	// WHILE
	if contentList.Len() == 1 || contentList.At(1).(objects.Token).Type != fract.TypeIn {
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
			if i.processCondition(&contentList) == grammar.TRUE {
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
