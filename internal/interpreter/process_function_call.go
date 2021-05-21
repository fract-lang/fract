package interpreter

import (
	"fmt"
	"strings"

	"github.com/fract-lang/fract/internal/functions/embed"
	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// isParamSet Argument type is param set?
func isParamSet(tokens []objects.Token) bool {
	return tokens[0].Type == fract.TypeName && tokens[1].Value == grammar.TokenEquals
}

// getParamsArgumentValues decompose and returns params values.
func (i *Interpreter) getParamsArgumentValues(tokens []objects.Token, index, braceCount, lastComma *int) objects.Value {
	returnValue := objects.Value{
		Content: []objects.DataFrame{},
		Array:   true,
	}

	for ; *index < len(tokens); *index++ {
		current := tokens[*index]
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLParenthes ||
				current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket {
				*braceCount++
			} else {
				*braceCount--
			}
		} else if current.Type == fract.TypeComma && *braceCount == 0 {
			valueList := vector.Sublist(tokens, *lastComma, *index-*lastComma)
			if isParamSet(*valueList) {
				*index -= 4
				return returnValue
			}
			returnValue.Content = append(returnValue.Content, i.processValue(*valueList).Content...)
			*lastComma = *index + 1
		}
	}

	if *lastComma < len(tokens) {
		valueSlice := tokens[*lastComma:]
		if isParamSet(valueSlice) {
			*index -= 4
			return returnValue
		}
		returnValue.Content = append(returnValue.Content, i.processValue(valueSlice).Content...)
	}

	return returnValue
}

func (i *Interpreter) processArgument(function objects.Function, names *[]string, tokens []objects.Token,
	current objects.Token, index, count, braceCount, lastComma *int) objects.Variable {
	var paramSet bool

	length := *index - *lastComma
	if length < 1 {
		fract.Error(current, "Value is not defined!")
	} else if *count >= len(function.Parameters) {
		fract.Error(current, "Argument overflow!")
	}

	parameter := function.Parameters[*count]
	variable := objects.Variable{Name: parameter.Name}
	valueList := *vector.Sublist(tokens, *lastComma, length)
	current = valueList[0]

	// Check param set.
	if length >= 2 && isParamSet(valueList) {
		length -= 2
		if length < 1 {
			fract.Error(current, "Value is not defined!")
		}

		for _, parameter := range function.Parameters {
			if parameter.Name == current.Value {
				for _, name := range *names {
					if name == current.Value {
						fract.Error(current, "Keyword argument repeated!")
					}
				}
				*count++
				paramSet = true
				*names = append(*names, current.Value)
				returnValue := objects.Variable{Name: current.Value}
				//Parameter is params typed?
				if parameter.Params {
					*lastComma += 2
					returnValue.Value = i.getParamsArgumentValues(tokens, index, braceCount, lastComma)
				} else {
					returnValue.Value = i.processValue(valueList[2:])
				}
				return returnValue
			}
		}

		fract.Error(current, "Parameter is not defined in this name!: "+current.Value)
	}

	if paramSet {
		fract.Error(current, "After the parameter has been given a special value, all parameters must be shown privately!")
	}

	*count++
	*names = append(*names, variable.Name)
	// Parameter is params typed?
	if parameter.Params {
		variable.Value = i.getParamsArgumentValues(tokens, index, braceCount, lastComma)
	} else {
		variable.Value = i.processValue(valueList)
	}
	return variable
}

// processFunctionCall call function and returns returned value.
func (i *Interpreter) processFunctionCall(tokens []objects.Token) objects.Value {
	_name := tokens[0]

	// Name is not defined?
	nameIndex, source := i.functionIndexByName(_name)
	if nameIndex == -1 {
		fract.Error(_name, "Function is not defined in this name!: "+_name.Value)
	}

	function := source.functions[nameIndex]

	var (
		vars  []objects.Variable
		names = new([]string)
		count = new(int)
	)

	// Decompose arguments.
	if tokens, _ = parser.DecomposeBrace(&tokens, grammar.TokenLParenthes, grammar.TokenRParenthes, false); tokens != nil {
		var (
			braceCount = new(int)
			lastComma  = new(int)
		)

		for index := 0; index < len(tokens); index++ {
			current := tokens[index]
			if current.Type == fract.TypeBrace {
				if current.Value == grammar.TokenLParenthes ||
					current.Value == grammar.TokenLBrace ||
					current.Value == grammar.TokenLBracket {
					*braceCount++
				} else {
					*braceCount--
				}
			} else if current.Type == fract.TypeComma && *braceCount == 0 {
				vars = append(vars, i.processArgument(function, names, tokens, current, &index, count, braceCount, lastComma))
				*lastComma = index + 1
			}
		}

		if *lastComma < len(tokens) {
			tokenLen := len(tokens)
			vars = append(vars, i.processArgument(function, names, tokens, tokens[*lastComma], &tokenLen, count, braceCount, lastComma))
		}
	}

	// All parameters is not defined?
	if *count < len(function.Parameters)-function.DefaultParameterCount {
		var sb strings.Builder
		sb.WriteString("All required positional parameters is not defined:")
		for _, parameter := range function.Parameters {
			if parameter.Default.Content != nil {
				break
			}
			argMsg := " '" + parameter.Name + "',"
			for _, name := range *names {
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
	for ; *count < len(function.Parameters); *count++ {
		current := function.Parameters[*count]
		if current.Default.Content != nil {
			vars = append(vars,
				objects.Variable{
					Name:  current.Name,
					Value: current.Default,
				})
		}
	}

	returnValue := objects.Value{}

	// Is embed function?
	if function.Tokens == nil {
		// Add name token for exceptions.
		function.Tokens = [][]objects.Token{{_name}}

		switch source.Lexer.File.Path {
		default: //* Direct embed functions.
			switch function.Name {
			case "print":
				embed.Print(function, vars)
			case "input":
				returnValue = embed.Input(function, vars)
			case "len":
				returnValue = embed.Len(function, vars)
			case "range":
				returnValue = embed.Range(function, vars)
			case "make":
				returnValue = embed.Make(function, vars)
			case "string":
				returnValue = embed.String(function, vars)
			case "int":
				returnValue = embed.Int(function, vars)
			case "float":
				returnValue = embed.Float(function, vars)
			default:
				embed.Exit(function, vars)
			}
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
		block := except.Block{
			Try: func() {
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
			},
		}
		block.Do()

		source.Tokens = itokens

		// Remove temporary functions.
		source.functions = source.functions[:functionLen]

		// Remove temporary variables.
		source.variables = variables

		source.functionCount--
		source.funcTempVariables = old
		source.index = nameIndex

		if block.Exception != nil {
			panic(fmt.Errorf(block.Exception.Message))
		}
	}

	return returnValue
}
