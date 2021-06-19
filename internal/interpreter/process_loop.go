package interpreter

import (
	"fmt"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// processKwState returns return value of kwstate.
func processKwState(kwstate uint8) uint8 {
	if kwstate != fract.FUNCReturn {
		return fract.TypeNone
	}
	return kwstate
}

// processLoop process loop blocks and returns keyword state.
func (i *Interpreter) processLoop(tokens []objects.Token) uint8 {
	// Content is empty?
	if vtokens := vector.Sublist(tokens, 1, len(tokens)-1); vtokens == nil {
		tokens = nil
	} else {
		tokens = *vtokens
	}
	functionLen := len(i.functions)
	_break := false
	kwstate := fract.TypeNone
	iindex := i.index

	//*************
	//    WHILE
	//*************
	if tokens == nil || len(tokens) >= 1 {
		if tokens == nil || len(tokens) == 1 || len(tokens) >= 1 && tokens[1].Type != fract.TypeIn && tokens[1].Type != fract.TypeComma {
			variableLen := len(i.variables)
			/* Infinity loop. */
			if tokens == nil {
				for {
					i.index++
					tokens := i.Tokens[i.index]
					if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
						// Remove temporary variables.
						i.variables = i.variables[:variableLen]
						// Remove temporary functions.
						i.functions = i.functions[:functionLen]
						if _break {
							return processKwState(kwstate)
						}
						i.index = iindex
						continue
					} else if tokens[0].Type == fract.TypeElse { // Else block.
						if len(tokens) > 1 {
							fract.Error(tokens[0], "Else block is not take any arguments!")
						}
						i.skipBlock(false)
						i.index--
						continue
					}
					kwstate = i.processTokens(tokens)
					if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
						_break = true
						i.skipBlock(false)
						i.index--
					} else if kwstate == fract.LOOPContinue { // Continue loop?
						i.skipBlock(false)
						i.index--
					}
				}
			}

			/* Interpret/skip block. */
			conditionList := tokens
			condition := i.processCondition(conditionList)
			_else := condition == "false"
			for {
				i.index++
				tokens := i.Tokens[i.index]

				if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
					// Remove temporary variables.
					i.variables = i.variables[:variableLen]
					// Remove temporary functions.
					i.functions = i.functions[:functionLen]
					condition = i.processCondition(conditionList)
					if _break || condition != "true" {
						return processKwState(kwstate)
					}
					i.index = iindex
					continue
				} else if tokens[0].Type == fract.TypeElse { // Else block.
					if len(tokens) > 1 {
						fract.Error(tokens[0], "Else block is not take any arguments!")
					}
					if condition == "true" {
						i.skipBlock(false)
						i.index--
						continue
					}
					// Remove temporary variables.
					i.variables = i.variables[:variableLen]
					// Remove temporary functions.
					i.functions = i.functions[:functionLen]
					if !_else {
						i.skipBlock(false)
						return processKwState(kwstate)
					}
					for {
						i.index++
						tokens = i.Tokens[i.index]
						if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
							// Remove temporary variables.
							i.variables = i.variables[:variableLen]
							// Remove temporary functions.
							i.functions = i.functions[:functionLen]
							return processKwState(kwstate)
						}
						kwstate = i.processTokens(tokens)
						if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
							_break = true
							i.skipBlock(false)
							i.index--
						} else if kwstate == fract.LOOPContinue { // Continue loop?
							i.skipBlock(false)
							i.index--
						}
					}
				}

				// Condition is true?
				if condition == "true" {
					kwstate = i.processTokens(tokens)
					if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
						_break = true
						i.skipBlock(false)
						i.index--
					} else if kwstate == fract.LOOPContinue { // Continue loop?
						i.skipBlock(false)
						i.index--
					}
				} else {
					if _else {
						i.skipBlock(true)
						continue
					}
					_break = true
					i.skipBlock(false)
					i.index--
				}
			}
		}
	}

	//*************
	//   FOREACH
	//*************
	nameToken := tokens[0]
	// Name is not name?
	if nameToken.Type != fract.TypeName {
		fract.Error(nameToken, "This is not a valid name!")
	}
	// Element name?
	elementName := ""
	if tokens[1].Type == fract.TypeComma {
		if len(tokens) < 3 || tokens[2].Type != fract.TypeName {
			fract.Error(tokens[1], "Element name is not defined!")
		}
		if tokens[2].Value != "_" {
			elementName = tokens[2].Value
		}
		if len(tokens)-3 == 0 {
			tokens[2].Column += len(tokens[2].Value)
			fract.Error(tokens[2], "Value is not defined!")
		}
		tokens = tokens[2:]
	}
	if vtokens, inToken := vector.Sublist(tokens, 2, len(tokens)-2), tokens[1]; vtokens != nil {
		tokens = *vtokens
	} else {
		fract.Error(inToken, "Value is not defined!")
	}
	value := i.processValue(tokens)
	// Type is not array?
	if !value.Array && value.Content[0].Type != objects.VALString {
		fract.Error(tokens[0], "Foreach loop must defined array value!")
	}
	// Empty array?
	if value.Array && len(value.Content) == 0 ||
		value.Content[0].Type == objects.VALString && value.Content[0].Data == "" {
		varLen := len(i.variables)
		for {
			i.index++
			tokens := i.Tokens[i.index]
			if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
				return kwstate
			} else if tokens[0].Type == fract.TypeElse { // Else block.
				if len(tokens) > 1 {
					fract.Error(tokens[0], "Else block is not take any arguments!")
				}
				for {
					i.index++
					tokens = i.Tokens[i.index]
					if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
						// Remove temporary variables.
						i.variables = i.variables[:varLen]
						// Remove temporary functions.
						i.functions = i.functions[:functionLen]
						return processKwState(kwstate)
					}
					kwstate = i.processTokens(tokens)
					if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
						_break = true
						i.skipBlock(false)
						i.index--
					} else if kwstate == fract.LOOPContinue { // Continue loop?
						i.skipBlock(false)
						i.index--
					}
				}
			}
			i.skipBlock(true)
		}
	}

	i.variables = append(
		[]objects.Variable{
			{ // Index.
				Name: nameToken.Value,
				Value: objects.Value{
					Content: []objects.Data{{
						Data: "0",
						Type: objects.VALInteger,
					}},
				},
			},
			{ // Element.
				Name:  elementName,
				Value: objects.Value{},
			}}, i.variables...)
	vars := i.variables
	index := &i.variables[0]
	element := &i.variables[1]
	if index.Name == "_" {
		index.Name = ""
	}
	var length int
	if value.Array {
		length = len(value.Content)
	} else {
		length = len(value.Content[0].String())
	}
	if element.Name != "" {
		if value.Array {
			element.Value.Content = []objects.Data{value.Content[0]}
		} else {
			element.Value.Content = []objects.Data{{
				Data: string(value.Content[0].String()[0]),
				Type: objects.VALString,
			}}
		}
	}
	//? Interpret block.
	for vindex := 0; vindex < length; {
		i.index++
		tokens := i.Tokens[i.index]
		if tokens[0].Type == fract.TypeBlockEnd { // Block is ended.
			// Remove temporary variables.
			i.variables = vars
			// Remove temporary functions.
			i.functions = i.functions[:functionLen]
			vindex++
			if _break ||
				(value.Array && vindex == len(value.Content) ||
					!value.Array && vindex == len(value.Content[0].String())) {
				break
			}
			i.index = iindex
			if index.Name != "" {
				index.Value.Content = []objects.Data{{
					Data: fmt.Sprint(vindex),
					Type: objects.VALInteger,
				}}
			}
			if element.Name != "" {
				if value.Array {
					element.Value.Content = []objects.Data{value.Content[vindex]}
				} else {
					element.Value.Content = []objects.Data{{
						Data: string(value.Content[0].String()[vindex]),
						Type: objects.VALString,
					}}
				}
			}
			continue
		} else if tokens[0].Type == fract.TypeElse { // Else block.
			if len(tokens) > 1 {
				fract.Error(tokens[0], "Else block is not take any arguments!")
			}
			i.skipBlock(false)
			i.index--
			continue
		}
		kwstate = i.processTokens(tokens)
		if kwstate == fract.LOOPBreak || kwstate == fract.FUNCReturn { // Break loop or return?
			_break = true
			i.skipBlock(false)
			i.index--
		} else if kwstate == fract.LOOPContinue { // Continue next?
			i.skipBlock(false)
			i.index--
		}
	}
	// Remove loop variables.
	i.variables = vars[2:]
	return processKwState(kwstate)
}
