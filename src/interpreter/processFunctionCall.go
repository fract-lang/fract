/*
	processFunctionCall Function.
*/

package interpreter

import (
	"../fract"
	"../fract/name"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processFunctionCall Process function call.
// tokens Tokens to process.
func (i *Interpreter) processFunctionCall(tokens *vector.Vector) {
	_name := tokens.Vals[0].(objects.Token)

	// Name is not defined?
	nameIndex := name.FunctionIndexByName(i.funcs, _name.Value)
	if nameIndex == -1 {
		fract.Error(_name, "Function is not defined!: "+_name.Value)
	}

	tokens, _ = parser.DecomposeBrace(tokens, grammar.TokenLParenthes,
		grammar.TokenRParenthes, false)

	i.blockCount++

	function := i.funcs.Vals[nameIndex].(objects.Function)
	variableLen := len(i.vars.Vals)
	functionLen := len(i.funcs.Vals)

	nameIndex = i.index
	i.index = function.Start

	for ; i.index < len(i.tokens.Vals); i.index++ {
		tokens = i.tokens.Vals[i.index].(*vector.Vector)

		first := tokens.Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			// Remove temporary variables.
			i.vars.Vals = i.vars.Vals[:variableLen]
			// Remove temporary functions.
			i.funcs.Vals = i.funcs.Vals[:functionLen]

			i.subtractBlock(&first)
			break
		}

		i.processTokens(tokens, true)
	}

	i.index = nameIndex
}
