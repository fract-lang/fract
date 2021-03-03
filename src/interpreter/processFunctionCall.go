/*
	processFunctionCall Function.
*/

package interpreter

import (
	"../fract"
	"../fract/dt"
	"../fract/name"
	"../grammar"
	"../objects"
	"../parser"
	"../utilities/vector"
)

// processFunctionCall Process function call.
// tokens Tokens to process.
func (i *Interpreter) processFunctionCall(tokens *vector.Vector) objects.Value {
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
	returnValue := objects.Value{
		Content: nil,
	}

	nameIndex = i.index
	i.index = function.Start

	for ; i.index < len(i.tokens.Vals); i.index++ {
		tokens = i.tokens.Vals[i.index].(*vector.Vector)

		first := tokens.Vals[0].(objects.Token)
		if first.Type == fract.TypeBlockEnd { // Block is ended.
			break
		}

		if i.processTokens(tokens, true) == fract.FUNCReturn {
			valueList := tokens.Sublist(1, len(tokens.Vals)-1)
			if len(valueList.Vals) == 0 {
				break
			}
			returnValue = i.processValue(valueList)
			// Check value and data type compatibility.
			if dt.IsIntegerType(function.ReturnType) && returnValue.Type != fract.VTInteger {
				fract.Error(first, "Return data type and value data type is not compatible!")
			}
			break
		}
	}

	// Remove temporary variables.
	i.vars.Vals = i.vars.Vals[:variableLen]
	// Remove temporary functions.
	i.funcs.Vals = i.funcs.Vals[:functionLen]

	i.subtractBlock(nil)

	if function.ReturnType != "" && returnValue.Content == nil {
		fract.Error(
			i.tokens.Vals[function.Start-1].(*vector.Vector).Vals[0].(objects.Token),
			"This function is returnable but not return anything!")
	}

	i.index = nameIndex
	return returnValue
}
