/*
	processTryCatch Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
)

// processTryCatch Process Try-Catch block.
// tokens Tokens to process.
func (i *Interpreter) processTryCatch(tokens []obj.Token) int16 {
	if len(tokens) > 1 {
		fract.Error(tokens[1], "Invalid syntax!")
	}

	fract.TryCount++

	variableLen := len(i.variables)
	functionLen := len(i.functions)
	kwstate := fract.TypeNone

	except.Block{
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
		Catch: func(e obj.Exception) {
			i.loopCount = 0
			fract.TryCount--
			i.variables = i.variables[:variableLen]
			i.functions = i.functions[:functionLen]

			count := 0
			for {
				i.index++
				tokens := i.Tokens[i.index]
				if tokens[0].Type == fract.TypeBlockEnd {
					count--
					if count < 0 {
						break
					}
				} else if parser.IsBlockStatement(tokens) {
					count++
				}

				if count > 0 {
					continue
				}

				if tokens[0].Type == fract.TypeCatch {
					i.index--
					break
				}
			}

			if count < 0 {
				return
			}

			for {
				i.index++
				tokens := i.Tokens[i.index]

				if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
					break
				} else if tokens[0].Type == fract.TypeCatch { // Catch.
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

					break
				}
			}

			i.variables = i.variables[:variableLen]
			i.functions = i.functions[:functionLen]
		},
	}.Do()

	return kwstate
}
