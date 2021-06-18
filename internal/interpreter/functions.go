package interpreter

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/fract-lang/fract/internal/functions/embed"
	"github.com/fract-lang/fract/pkg/except"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// Instance for function calls.
type functionCall struct {
	function objects.Function
	name     objects.Token
	source   *Interpreter
	args     []objects.Variable
}

func (c functionCall) call() objects.Value {
	returnValue := objects.Value{}
	// Is embed function?
	if c.function.Tokens == nil {
		// Add name token for exceptions.
		c.function.Tokens = [][]objects.Token{{c.name}}
		switch c.function.Name {
		case "print":
			embed.Print(c.function, c.args)
		case "input":
			returnValue = embed.Input(c.function, c.args)
		case "len":
			returnValue = embed.Len(c.function, c.args)
		case "range":
			returnValue = embed.Range(c.function, c.args)
		case "make":
			returnValue = embed.Make(c.function, c.args)
		case "string":
			returnValue = embed.String(c.function, c.args)
		case "int":
			returnValue = embed.Int(c.function, c.args)
		case "float":
			returnValue = embed.Float(c.function, c.args)
		default:
			embed.Exit(c.function, c.args)
		}
	} else {
		// Process block.
		variables := c.source.variables
		deferLen := len(defers)
		if c.source.funcTempVariables == 0 {
			c.source.funcTempVariables = len(c.source.variables)
		}
		c.source.variables = append(c.args, c.source.variables[:c.source.funcTempVariables]...)
		c.source.functionCount++
		old := c.source.funcTempVariables
		c.source.funcTempVariables = len(c.args)
		functionLen := len(c.source.functions)
		nameIndex := c.source.index
		itokens := c.source.Tokens
		c.source.Tokens = c.function.Tokens
		c.source.index = -1
		// Interpret block.
		block := except.Block{
			Try: func() {
				for {
					c.source.index++
					tokens := c.source.Tokens[c.source.index]
					c.source.funcTempVariables = len(c.source.variables) - c.source.funcTempVariables
					/*for _, t := range tokens {
						fmt.Print(t.Value)
					}
					fmt.Println()*/
					if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
						break
					} else if c.source.processTokens(tokens) == fract.FUNCReturn {
						if c.source.returnValue == nil {
							break
						}
						returnValue = *c.source.returnValue
						c.source.returnValue = nil
						break
					}
				}
			},
		}
		block.Do()
		c.source.Tokens = itokens
		// Remove temporary functions.
		c.source.functions = c.source.functions[:functionLen]
		// Remove temporary variables.
		c.source.variables = variables
		c.source.functionCount--
		c.source.funcTempVariables = old
		c.source.index = nameIndex
		if block.Exception != nil {
			defers = defers[:deferLen]
			panic(fmt.Errorf(block.Exception.Message))
		}
		for index := len(defers) - 1; index >= deferLen; index-- {
			defers[index].call()
		}
		defers = defers[:deferLen]
	}
	return returnValue
}

// isParamSet Argument type is param set?
func isParamSet(tokens []objects.Token) bool {
	return tokens[0].Type == fract.TypeName && tokens[1].Value == "="
}

// getParamsArgumentValues decompose and returns params values.
func (i *Interpreter) getParamsArgumentValues(tokens []objects.Token, index, braceCount, lastComma *int) objects.Value {
	returnValue := objects.Value{
		Content: []objects.Data{},
		Array:   true,
	}
	for ; *index < len(tokens); *index++ {
		current := tokens[*index]
		if current.Type == fract.TypeBrace {
			if current.Value == "(" || current.Value == "{" || current.Value == "[" {
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
			val := i.processValue(*valueList)
			if val.Array {
				returnValue.Content = append(returnValue.Content, objects.Data{
					Data: val.Content,
					Type: objects.VALArray,
				})
			} else {
				returnValue.Content = append(returnValue.Content, val.Content...)
			}
			*lastComma = *index + 1
		}
	}
	if *lastComma < len(tokens) {
		valueSlice := tokens[*lastComma:]
		if isParamSet(valueSlice) {
			*index -= 4
			return returnValue
		}
		val := i.processValue(valueSlice)
		if val.Array {
			returnValue.Content = append(returnValue.Content, objects.Data{
				Data: val.Content,
				Type: objects.VALArray,
			})
		} else {
			returnValue.Content = append(returnValue.Content, val.Content...)
		}
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
		fract.Error(current, "Parameter is not defined in this name: "+current.Value)
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

// Process function call model and initialize moden instance.
func (i *Interpreter) processFunctionCallModel(tokens []objects.Token) functionCall {
	_name := tokens[0]
	// Name is not defined?
	nameIndex, source := i.functionIndexByName(_name)
	var function objects.Function
	if nameIndex == -1 {
		_name := _name
		if index := strings.Index(_name.Value, "."); index != -1 {
			if i.importIndexByName(_name.Value[:index]) == -1 {
				fract.Error(_name, "'"+_name.Value[:index]+"' is not defined!")
			}
			source = i.Imports[i.importIndexByName(_name.Value[:index])].Source
			_name.Value = _name.Value[index+1:]
			for _, current := range source.variables {
				if unicode.IsUpper(rune(current.Name[0])) && current.Name == _name.Value && !current.Value.Array && current.Value.Content[0].Type == objects.VALFunction {
					_name.File = nil
					function = current.Value.Content[0].Data.(objects.Function)
					break
				}
			}
		} else {
			for _, current := range i.variables {
				if current.Name == _name.Value && !current.Value.Array && current.Value.Content[0].Type == objects.VALFunction {
					_name.File = nil
					function = current.Value.Content[0].Data.(objects.Function)
					source = i
					break
				}
			}
		}
		if _name.File != nil {
			fract.Error(_name, "Function is not defined in this name: "+_name.Value)
		}
	} else {
		function = source.functions[nameIndex]
	}
	var (
		names = new([]string)
		count = new(int)
		args  []objects.Variable
	)
	// Decompose arguments.
	if tokens, _ = parser.DecomposeBrace(&tokens, "(", ")", false); tokens != nil {
		var (
			braceCount = new(int)
			lastComma  = new(int)
		)
		for index := 0; index < len(tokens); index++ {
			current := tokens[index]
			if current.Type == fract.TypeBrace {
				if current.Value == "(" || current.Value == "{" || current.Value == "[" {
					*braceCount++
				} else {
					*braceCount--
				}
			} else if current.Type == fract.TypeComma && *braceCount == 0 {
				args = append(args, i.processArgument(function, names, tokens, current, &index, count, braceCount, lastComma))
				*lastComma = index + 1
			}
		}
		if *lastComma < len(tokens) {
			tokenLen := len(tokens)
			args = append(args, i.processArgument(function, names, tokens, tokens[*lastComma], &tokenLen, count, braceCount, lastComma))
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
			args = append(args,
				objects.Variable{
					Name:  current.Name,
					Value: current.Default,
				})
		}
	}

	return functionCall{
		function: function,
		name:     _name,
		source:   source,
		args:     args,
	}
}

// processFunctionCall call function and returns returned value.
func (i *Interpreter) processFunctionCall(tokens []objects.Token) objects.Value {
	return i.processFunctionCallModel(tokens).call()
}

func (i *Interpreter) processFunctionDeclaration(tokens []objects.Token, protected bool) {
	tokenLen := len(tokens)
	_name := tokens[1]
	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	} else if strings.Contains(_name.Value, ".") {
		fract.Error(_name, "Names is cannot include dot!")
	}
	// Name is already defined?
	if line := i.definedName(_name); line != -1 {
		fract.Error(_name, "\""+_name.Value+"\" is already defined at line: "+fmt.Sprint(line))
	}
	// Function parentheses are not defined?
	if tokenLen < 4 {
		fract.Error(_name, "Where is the function parentheses?")
	}
	i.index++
	function := objects.Function{
		Name:       _name.Value,
		Line:       i.index,
		Parameters: []objects.Parameter{},
		Protected:  protected,
	}
	dtToken := tokens[tokenLen-1]
	if dtToken.Type != fract.TypeBrace || dtToken.Value != ")" {
		fract.Error(dtToken, "Invalid syntax!")
	}
	if paramList := vector.Sublist(tokens, 3, tokenLen-4); paramList != nil {
		paramList := *paramList
		// Decompose function parameters.
		paramName, defaultDefined := true, false
		var lastParameter objects.Parameter
		for index := 0; index < len(paramList); index++ {
			current := paramList[index]
			if paramName {
				if current.Type == fract.TypeParams {
					continue
				} else if current.Type != fract.TypeName {
					fract.Error(current, "Parameter name is not found!")
				}
				lastParameter = objects.Parameter{
					Name:   current.Value,
					Params: index > 0 && paramList[index-1].Type == fract.TypeParams,
				}
				function.Parameters = append(function.Parameters, lastParameter)
				paramName = false
				continue
			} else {
				paramName = true
				// Default value definition?
				if current.Value == "=" {
					brace := 0
					index++
					start := index
					for ; index < len(paramList); index++ {
						current = paramList[index]
						if current.Type == fract.TypeBrace {
							if current.Value == "{" || current.Value == "(" || current.Value == "[" {
								brace++
							} else {
								brace--
							}
						} else if current.Type == fract.TypeComma {
							break
						}
					}
					if index-start < 1 {
						fract.Error(paramList[start-1], "Value is not defined!")
					}
					lastParameter.Default = i.processValue(*vector.Sublist(paramList, start, index-start))
					if lastParameter.Params && !lastParameter.Default.Array {
						fract.Error(current, "Params parameter is can only take array values!")
					}
					function.Parameters[len(function.Parameters)-1] = lastParameter
					function.DefaultParameterCount++
					defaultDefined = true
					continue
				}
				if lastParameter.Default.Content == nil && defaultDefined {
					fract.Error(current, "All parameters after a given parameter with a default value must take a default value!")
				} else if current.Type != fract.TypeComma {
					fract.Error(current, "Comma is not found!")
				}
			}
		}
		if lastParameter.Default.Content == nil && defaultDefined {
			fract.Error(tokens[len(tokens)-1], "All parameters after a given parameter with a default value must take a default value!")
		}
	}
	i.skipBlock(false)
	function.Tokens = i.Tokens[function.Line : i.index+1]
	function.Line = _name.Line
	i.functions = append(i.functions, function)
}
