package interpreter

import (
	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
)

// processTryCatch process try-catch blocks and returns keyword state.
func (i *Interpreter) processTryCatch(tokens []objects.Token) uint8 {
	if len(tokens) > 1 {
		fract.Error(tokens[1], "Invalid syntax!")
	}

	fract.TryCount++

	variableLen := len(i.variables)
	functionLen := len(i.functions)
	kwstate := fract.TypeNone

	(&except.Block{
		Try: func() {
			for {
				i.index++
				tokens := i.Tokens[i.index]

				if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
					break
				} else if tokens[0].Type == fract.TypeCatch { // Catch.
					i.skipBlock(false)
					break
				}

				if kwstate = i.processTokens(tokens); kwstate != fract.TypeNone {
					i.skipBlock(false)
				}
			}

			fract.TryCount--
			i.variables = i.variables[:variableLen]
			i.functions = i.functions[:functionLen]
		},
		Catch: func(e *objects.Exception) {
			i.loopCount = 0
			fract.TryCount--
			i.variables = i.variables[:variableLen]
			i.functions = i.functions[:functionLen]

			count := 0
			for {
				i.index++
				tokens = i.Tokens[i.index]
				if tokens[0].Type == fract.TypeBlockEnd {
					count--
					if count < 0 { break }
				} else if parser.IsBlockStatement(tokens) {
					count++
				}

				if count > 0 { continue }
				if tokens[0].Type == fract.TypeCatch { break }
			}

			// Ended block.
			if count < 0 { return }

			// Catch block.

			if len(tokens) > 1 {
				fract.Error(tokens[1], "Invalid syntax!")
			}

			for {
				i.index++
				tokens := i.Tokens[i.index]

				if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
					break
				}

				if kwstate = i.processTokens(tokens); kwstate != fract.TypeNone {
					i.skipBlock(false)
				}
			}

			i.variables = i.variables[:variableLen]
			i.functions = i.functions[:functionLen]
		},
	}).Do()

	return kwstate
}
