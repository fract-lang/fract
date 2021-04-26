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
	nameIndex, source := i.functionIndexByName(_name)
	if nameIndex == -1 {
		fract.Error(_name, "Function is not defined in this name!: "+_name.Value)
	}

	function := source.functions[nameIndex]

	var (
		vars  []*obj.Variable
		names []string
		count int
	)

	// Decompose arguments.
	if tokens, _ = parser.DecomposeBrace(&tokens, grammar.TokenLParenthes,
		grammar.TokenRParenthes, false); tokens != nil {
		braceCount, lastComma, tokenLen := 0, 0, len(tokens)
		paramSet := false

		// processArgument Process function argument.
		// current Current token.
		// count Count of appended arguments.
		// index Index of tokens state.
		processArgument := func(current obj.Token, index *int) *obj.Variable {
			getParamsArgumentValue := func() obj.Value {
				returnValue := obj.Value{
					Content: []obj.DataFrame{},
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

			if count >= len(function.Parameters) {
				fract.Error(current, "Argument overflow!")
			}

			parameter := function.Parameters[count]
			variable := &obj.Variable{Name: parameter.Name}
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
						count++
						valueList = valueList[2:]
						paramSet = true
						names = append(names, current.Value)
						returnValue := &obj.Variable{Name: current.Value}
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

			count++
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
	if count < len(function.Parameters)-function.DefaultParameterCount {
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
			vars = append(vars, &obj.Variable{
				Name:  current.Name,
				Value: current.Default,
			})
		}
	}

	returnValue := obj.Value{}

	// Is embed function?
	if function.Tokens == nil {
		// Add name token for exceptions.
		function.Tokens = [][]obj.Token{{_name}}

		switch function.Name {
		case "print":
			functions.Print(function, vars)
		case "input":
			returnValue = functions.Input(function, vars)
		case "len":
			returnValue = functions.Len(function, vars)
		case "range":
			returnValue = functions.Range(function, vars)
		case "make":
			returnValue = functions.Make(function, vars)
		case "string":
			returnValue = functions.String(function, vars)
		case "int":
			returnValue = functions.Int(function, vars)
		case "float":
			returnValue = functions.Float(function, vars)
		default:
			functions.Exit(function, vars)
		}
	} else {
		// Process block.

		variables := source.variables
		if source.funcTempVariables == 0 {
			source.funcTempVariables = len(source.variables)
		}
		source.variables = append(vars, source.variables[:source.funcTempVariables]...)

		source.functionCount++

		old := source.funcTempVariables
		source.funcTempVariables = len(vars)

		functionLen := len(source.functions)
		nameIndex = source.index
		itokens := source.Tokens
		source.Tokens = function.Tokens

		source.index = -1

		// Interpret block.
		for {
			source.index++
			tokens := source.Tokens[source.index]
			source.funcTempVariables = len(source.variables) - source.funcTempVariables

			if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
				break
			} else if source.processTokens(tokens) == fract.FUNCReturn {
				if source.returnValue == nil {
					break
				}
				returnValue = *source.returnValue
				source.returnValue = nil
				break
			}
		}

		source.Tokens = itokens

		// Remove temporary functions.
		source.functions = source.functions[:functionLen]

		// Remove temporary variables.
		source.variables = variables

		source.functionCount--
		source.funcTempVariables = old
		source.index = nameIndex
	}

	return returnValue
}
