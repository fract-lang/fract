/*
	processFunctionCall Function.
*/

package interpreter

import (
	"strings"

	"github.com/fract-lang/fract/internal/functions"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// processFunctionCall Process function call.
// tokens Tokens to process.
func (i *Interpreter) processFunctionCall(tokens []obj.Token) obj.Value {
	_name := tokens[0]

	// Name is not defined?
	nameIndex := i.functionIndexByName(_name.Value)
	if nameIndex == -1 {
		fract.Error(_name, "Function is not defined in this name!: "+_name.Value)
	}

	function := i.functions[nameIndex]
	vars, names := make([]obj.Variable, 0), []string{}
	count := 0

	// Decompose arguments.
	if tokens, _ = parser.DecomposeBrace(&tokens, grammar.TokenLParenthes,
		grammar.TokenRParenthes, false); tokens != nil {
		braceCount, lastComma, tokenLen := 0, 0, len(tokens)
		paramSet := false

		// processArgument Process function argument.
		// current Current token.
		// count Count of appended arguments.
		// index Index of tokens state.
		processArgument := func(current obj.Token, index *int) obj.Variable {
			getParamsArgumentValue := func() obj.Value {
				returnValue := obj.Value{
					Content: []string{},
					Array:   true,
				}

				for ; *index < len(tokens); *index++ {
					current := tokens[*index]
					if current.Type == fract.TypeBrace {
						if current.Value == grammar.TokenLParenthes ||
							current.Value == grammar.TokenLBrace ||
							current.Value == grammar.TokenLBracket {
							braceCount++
						} else {
							braceCount--
						}
					} else if current.Type == fract.TypeComma && braceCount == 0 {
						valueList := vector.Sublist(tokens, lastComma, *index-lastComma)
						if parser.IsParamSet(*valueList) {
							*index -= 4
							return returnValue
						}
						returnValue.Content = append(returnValue.Content,
							i.processValue(valueList).Content...)
						lastComma = *index + 1
					}
				}

				if lastComma < tokenLen {
					valueSlice := tokens[lastComma:]
					if parser.IsParamSet(valueSlice) {
						*index -= 4
						return returnValue
					}
					returnValue.Content = append(returnValue.Content,
						i.processValue(&valueSlice).Content...)
				}

				return returnValue
			}

			length := *index - lastComma
			if length < 1 {
				fract.Error(current, "Value is not defined!")
			}

			if count > len(function.Parameters)-function.DefaultParameterCount {
				fract.Error(current, "Argument overflow!")
			}

			parameter := function.Parameters[count]
			variable := obj.Variable{Name: parameter.Name}
			valueList := *vector.Sublist(tokens, lastComma, length)
			current = valueList[0]

			// Check param set.
			if length >= 2 && parser.IsParamSet(valueList) {
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
						valueList = valueList[2:]
						paramSet = true
						names = append(names, current.Value)
						returnValue := obj.Variable{Name: current.Value}
						//Parameter is params typed?
						if parameter.Params {
							lastComma += 2
							returnValue.Value = getParamsArgumentValue()
						} else {
							returnValue.Value = i.processValue(&valueList)
						}
						return returnValue
					}
				}

				fract.Error(current, "Parameter is not defined in this name!: "+current.Value)
			}

			if paramSet {
				fract.Error(current,
					"After the parameter has been given a special value, all parameters must be shown privately!")
			}

			if function.Parameters[count].Default.Content == nil {
				count++
			}
			names = append(names, variable.Name)
			// Parameter is params typed?
			if parameter.Params {
				variable.Value = getParamsArgumentValue()
			} else {
				variable.Value = i.processValue(&valueList)
			}
			return variable
		}

		for index := 0; index < len(tokens); index++ {
			current := tokens[index]
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

		if lastComma < tokenLen {
			vars = append(vars, processArgument(tokens[lastComma], &tokenLen))
		}
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
			vars = append(vars, obj.Variable{
				Name:  current.Name,
				Value: current.Default,
			})
		}
	}

	returnValue := obj.Value{}
	variables := append(make([]obj.Variable, 0), i.variables...)
	i.variables = append(i.variables[:i.funcTempVariables], vars...)

	// Is embed function?
	if function.Tokens == nil {
		parameters := make([]obj.Value, 0)

		// Set parameter defaults to normal values.
		for _, param := range function.Parameters {
			parameters = append(parameters,
				i.variables[i.varIndexByName(param.Name)].Value)
		}

		// Add name token for exceptions.
		function.Tokens = [][]obj.Token{{_name}}

		if function.Name == "print" {
			functions.Print(function, parameters)
		} else if function.Name == "input" {
			returnValue = functions.Input(function, parameters)
		} else if function.Name == "len" {
			returnValue = functions.Len(function, parameters)
		} else if function.Name == "range" {
			returnValue = functions.Range(function, parameters)
		} else if function.Name == "make" {
			returnValue = functions.Make(function, parameters)
		} else {
			functions.Exit(function, parameters)
		}
	} else {
		// Process block.

		i.functionCount++

		old := i.funcTempVariables
		i.funcTempVariables = len(vars)

		functionLen := len(i.functions)
		nameIndex = i.index
		itokens := i.Tokens
		i.Tokens = function.Tokens

		i.index = -1

		for {
			i.index++
			tokens := i.Tokens[i.index]
			i.funcTempVariables = len(i.variables) - i.funcTempVariables

			if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
				break
			} else if i.processTokens(tokens) == fract.FUNCReturn {
				tokens := i.Tokens[i.returnIndex]
				i.returnIndex = fract.TypeNone
				valueList := tokens[1:]
				if valueList == nil {
					break
				}
				returnValue = i.processValue(&valueList)
				break
			}
		}

		i.Tokens = itokens

		// Remove temporary functions.
		i.functions = i.functions[:functionLen]

		i.functionCount--
		i.funcTempVariables = old
		i.index = nameIndex
	}

	// Remove temporary variables.
	i.variables = variables

	return returnValue
}
