/*
	processFunctionCall Function.
*/

package interpreter

import (
	"strings"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// processFunctionCall Process function call.
// tokens Tokens to process.
func (i *Interpreter) processFunctionCall(tokens vector.Vector) objects.Value {
	_name := tokens.Vals[0].(objects.Token)

	// Name is not defined?
	nameIndex := i.functionIndexByName(_name.Value)
	if nameIndex == -1 {
		fract.Error(_name, "Function is not defined in this name!: "+_name.Value)
	}

	function := i.funcs[nameIndex]

	// Decompose arguments.
	tokens, _ = parser.DecomposeBrace(&tokens, grammar.TokenLParenthes,
		grammar.TokenRParenthes, false)
	braceCount, lastComma, count := 0, 0, 0
	vars, names := make([]objects.Variable, 0), []string{}
	paramSet := false

	// processArgument Process function argument.
	// current Current token.
	// count Count of appended arguments.
	// index Index of tokens state.
	processArgument := func(current objects.Token, index *int) objects.Variable {
		length := *index - lastComma
		if length < 1 {
			fract.Error(current, "Value is not defined!")
		}

		if count > len(function.Parameters)-function.DefaultParameterCount {
			fract.Error(current, "Argument overflow!")
		}

		variable := objects.Variable{Name: function.Parameters[count].Name}
		valueList := tokens.Sublist(lastComma, length)
		current = valueList.Vals[0].(objects.Token)

		// Check param set.
		if current.Type == fract.TypeName {
			if length >= 2 {
				second := valueList.Vals[1].(objects.Token)
				if second.Value == grammar.TokenEquals {
					length -= 2
					if length < 1 {
						fract.Error(current, "Value is not defined!")
					}

					for _, parameter := range function.Parameters {
						if parameter.Name == current.Value {
							for _, name := range names {
								if name == current.Value {
									fract.Error(current, "Keyword argument repeated!")
								}
							}
							if parameter.Default.Content == nil {
								count++
							}
							valueList.Vals = valueList.Vals[2:]
							paramSet = true
							names = append(names, current.Value)
							return objects.Variable{
								Name:  current.Value,
								Value: i.processValue(valueList),
							}
						}
					}

					fract.Error(current, "Parameter is not defined in this name!")
				}
			}
		}

		if paramSet {
			fract.Error(current,
				"After the parameter has been given a special value, all parameters must be shown privately!")
		}

		if function.Parameters[count].Default.Content == nil {
			count++
		}
		names = append(names, variable.Name)
		variable.Value = i.processValue(valueList)
		return variable
	}

	for index := 0; index < len(tokens.Vals); index++ {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLParenthes ||
				current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket {
				braceCount++
			} else {
				braceCount--
			}
		} else if current.Type == fract.TypeComma && braceCount == 0 {
			vars = append(vars, processArgument(current, &index))
			lastComma = index + 1
		}
	}

	if tokenLen := len(tokens.Vals); lastComma < tokenLen {
		current := tokens.Vals[lastComma].(objects.Token)
		vars = append(vars, processArgument(current, &tokenLen))
	}

	// All parameters is not defined?
	if count != len(function.Parameters)-function.DefaultParameterCount {
		var sb strings.Builder
		sb.WriteString("All required positional parameters is not defined:")
		for _, parameter := range function.Parameters {
			if parameter.Default.Content != nil {
				break
			}
			argMsg := " '" + parameter.Name + "',"
			for _, name := range names {
				if parameter.Name == name {
					argMsg = ""
					break
				}
			}
			sb.WriteString(argMsg)
		}
		fract.Error(_name, sb.String()[:sb.Len()-1])
	}

	// Check default values.
	for ; count < len(function.Parameters); count++ {
		current := function.Parameters[count]
		if current.Default.Content != nil {
			vars = append(vars, objects.Variable{
				Name:  current.Name,
				Value: current.Default,
			})
		}
	}

	old := i.funcTempVariables
	variables := append(make([]objects.Variable, 0), i.vars...)
	i.vars = append(i.vars[:i.funcTempVariables], vars...)
	i.funcTempVariables = len(vars)

	var returnValue objects.Value
	functionLen := len(i.funcs)
	nameIndex = i.index
	itokens := i.tokens
	i.tokens.Vals = function.Tokens

	// Process block.
	i.functions++
	i.index = -1
	for {
		i.index++
		tokens := i.tokens.Vals[i.index].(vector.Vector)
		i.funcTempVariables = len(i.vars) - i.funcTempVariables

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
	i.vars = variables
	// Remove temporary functions.
	i.funcs = i.funcs[:functionLen]

	i.functions--
	i.funcTempVariables = old
	i.index = nameIndex
	i.tokens = itokens

	return returnValue
}
