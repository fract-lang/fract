/*
	processFunctionCall Function.
*/

package interpreter

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/grammar"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/parser"
	"github.com/fract-lang/fract/src/utils/vector"
)

// processArgument Process function argument.
// function Function.
// current Current token.
// count Count of appended arguments.
// value Value instance.
func processArgument(function objects.Function, current objects.Token,
	count int, value objects.Value) objects.Variable {
	if count >= len(function.Parameters) {
		fract.Error(current, "Argument overflow!")
	}
	return objects.Variable{
		Name:  function.Parameters[count],
		Const: false,
		Value: value,
	}
}

// processFunctionCall Process function call.
// tokens Tokens to process.
func (i *Interpreter) processFunctionCall(tokens vector.Vector) objects.Value {
	_name := tokens.Vals[0].(objects.Token)

	// Name is not defined?
	nameIndex := i.functionIndexByName(_name.Value)
	if nameIndex == -1 {
		fract.Error(_name, "Function is not defined in this name!: "+_name.Value)
	}

	function := i.funcs.Vals[nameIndex].(objects.Function)

	// Decompose arguments.
	tokens, _ = parser.DecomposeBrace(&tokens, grammar.TokenLParenthes,
		grammar.TokenRParenthes, false)
	braceCount := 0
	lastComma := 0
	count := 0
	vars := make([]interface{}, 0)
	for index, current := range tokens.Vals {
		current := current.(objects.Token)
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLParenthes ||
				current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket {
				braceCount++
			} else if current.Value == grammar.TokenRParenthes ||
				current.Value == grammar.TokenRBrace ||
				current.Value == grammar.TokenRBracket {
				braceCount--
			}
		} else if current.Type == fract.TypeComma && braceCount == 0 {
			valueList := tokens.Sublist(lastComma, index-lastComma)
			if valueList.Vals == nil {
				fract.Error(current, "Value is not defined!")
			}
			vars = append(vars, processArgument(function, current, count,
				i.processValue(valueList)))
			count++
			lastComma = index + 1
		}
	}

	if tokenLen := len(tokens.Vals); lastComma < tokenLen {
		current := tokens.Vals[lastComma].(objects.Token)
		valueList := tokens.Sublist(lastComma, tokenLen-lastComma)
		if valueList.Vals == nil {
			fract.Error(current, "Value is not defined!")
		}
		vars = append(vars, processArgument(function, current, count,
			i.processValue(valueList)))
		count++
	}

	// All parameters is not defined?
	if count != len(function.Parameters) {
		fract.Error(_name, "All parameters is not defined!")
	}

	old := i.funcTempVariables
	variables := append(make([]interface{}, 0), i.vars.Vals...)
	i.vars.Vals = append(i.vars.Vals[:i.funcTempVariables], vars...)
	i.funcTempVariables = len(vars)

	functionLen := len(i.funcs.Vals)
	returnValue := objects.Value{
		Content: nil,
	}

	nameIndex = i.index
	itokens := i.tokens
	i.tokens.Vals = function.Tokens

	// Process block.
	i.functions++
	i.index = -1
	for {
		i.index++
		tokens := i.tokens.Vals[i.index].(vector.Vector)
		i.funcTempVariables = len(i.vars.Vals) - i.funcTempVariables

		if tokens.Vals[0].(objects.Token).Type == fract.TypeBlockEnd { // Block is ended.
			break
		} else if i.processTokens(tokens) == fract.FUNCReturn {
			tokens := i.tokens.Vals[i.returnIndex].(vector.Vector)
			i.returnIndex = fract.TypeNone
			valueList := vector.Vector{Vals: tokens.Vals[1:]}
			if valueList.Vals == nil {
				break
			}
			returnValue = i.processValue(&valueList)
			break
		}
	}

	// Remove temporary variables.
	i.vars.Vals = variables
	// Remove temporary functions.
	i.funcs.Vals = i.funcs.Vals[:functionLen]

	i.functions--
	i.funcTempVariables = old
	i.index = nameIndex
	i.tokens = itokens

	return returnValue
}
