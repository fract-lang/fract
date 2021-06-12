package interpreter

import (
	"fmt"
	"runtime"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

func (i *Interpreter) processMacroIf(tokens []objects.Token) uint8 {
	tokenLen := len(tokens)
	conditionList := vector.Sublist(tokens, 1, tokenLen-1)
	// Condition is empty?
	if conditionList == nil {
		first := tokens[0]
		fract.ErrorCustom(first.File, first.Line, first.Column+len(first.Value), "Condition is empty!")
	}
	variables := i.variables
	functions := i.functions
	i.variables = append([]objects.Variable{
		{
			Name: "OS",
			Value: objects.Value{
				Content: []objects.Data{
					{
						Data: runtime.GOOS,
						Type: objects.VALString,
					},
				},
			},
		},
		{
			Name: "ARCH",
			Value: objects.Value{
				Content: []objects.Data{
					{
						Data: runtime.GOARCH,
						Type: objects.VALString,
					},
				},
			},
		},
	}, i.macroDefines...)
	state := i.processCondition(*conditionList)
	kwstate := fract.TypeNone
	/* Interpret/skip block. */
	for {
		i.index++
		tokens := i.Tokens[i.index]
		first := tokens[0]
		if first.Type == fract.TypeMacro {
			tokens := tokens[1:]
			first = tokens[0]
			if first.Type == fract.TypeBlockEnd { // Block is ended.
				goto ret
			} else if first.Type == fract.TypeElseIf { // Else if block.
				tokenLen = len(tokens)
				conditionList := vector.Sublist(tokens, 1, tokenLen-1)
				// Condition is empty?
				if conditionList == nil {
					first := tokens[0]
					fract.ErrorCustom(first.File, first.Line,
						first.Column+len(first.Value), "Condition is empty!")
				}
				if state == grammar.KwTrue {
					i.skipBlock(false)
					goto ret
				}
				state = i.processCondition(*conditionList)
				// Interpret/skip block.
				for {
					i.index++
					tokens := i.Tokens[i.index]
					first := tokens[0]
					if first.Type == fract.TypeMacro {
						tokens := tokens[1:]
						first = tokens[0]
						if first.Type == fract.TypeBlockEnd { // Block is ended.
							goto ret
						} else if first.Type == fract.TypeIf { // If block.
							if state == grammar.KwTrue && kwstate == fract.TypeNone {
								i.processMacroIf(tokens)
							} else {
								i.skipBlock(true)
							}
							continue
						} else if first.Type == fract.TypeElseIf || first.Type == fract.TypeElse { // Else if or else block.
							i.index--
							break
						}
					}
					// Condition is true?
					if state == grammar.KwTrue && kwstate == fract.TypeNone {
						i.variables, variables = variables, i.variables
						kwstate = i.processTokens(tokens)
						i.variables, variables = variables, i.variables
						if kwstate != fract.TypeNone {
							i.skipBlock(false)
						}
					} else {
						i.skipBlock(true)
					}
				}
				if state == grammar.KwTrue {
					i.skipBlock(false)
					goto ret
				}
				continue
			} else if first.Type == fract.TypeElse { // Else block.
				if len(tokens) > 1 {
					fract.Error(first, "Else block is not take any arguments!")
				}
				if state == grammar.KwTrue {
					i.skipBlock(false)
					goto ret
				}
				/* Interpret/skip block. */
				for {
					i.index++
					tokens := i.Tokens[i.index]
					first := tokens[0]
					if first.Type == fract.TypeMacro {
						tokens = tokens[1:]
						first = tokens[0]
						if first.Type == fract.TypeBlockEnd { // Block is ended.
							goto ret
						} else if first.Type == fract.TypeIf { // If block.
							if kwstate == fract.TypeNone {
								i.processMacroIf(tokens)
							} else {
								i.skipBlock(true)
							}
							continue
						}
					}
					// Condition is true?
					if kwstate == fract.TypeNone {
						i.variables, variables = variables, i.variables
						kwstate = i.processTokens(tokens)
						i.variables, variables = variables, i.variables
						if kwstate != fract.TypeNone {
							i.skipBlock(false)
						}
					}
				}
			}
		}
		// Condition is true?
		if state == grammar.KwTrue && kwstate == fract.TypeNone {
			i.variables, variables = variables, i.variables
			kwstate = i.processTokens(tokens)
			i.variables, variables = variables, i.variables
			if kwstate != fract.TypeNone {
				i.skipBlock(false)
			}
		} else {
			i.skipBlock(true)
		}
	}
ret:
	i.variables = variables
	i.functions = functions
	return kwstate
}

func (i *Interpreter) processMacroDefine(tokens []objects.Token) objects.Variable {
	if len(tokens) < 2 {
		fract.Error(tokens[0], "Define name is not defined!")
	}
	name := tokens[1]
	if name.Type != fract.TypeName {
		fract.Error(name, "Invalid name!")
	}
	// Exists name.
	for _, macro := range i.macroDefines {
		if macro.Name == name.Value {
			fract.Error(name, "This macro define is already defined in this name at line: "+fmt.Sprint(macro.Line))
		}
	}
	macro := objects.Variable{
		Name: name.Value,
		Line: name.Line,
	}
	if len(tokens) > 2 {
		variables := i.variables
		macro.Value = i.processValue(tokens[2:])
		i.variables = variables
	} else {
		macro.Value.Content = []objects.Data{
			{
				Data: grammar.KwFalse,
				Type: objects.VALBoolean,
			},
		}
	}
	return macro
}

// processMacro process macros and returns keyword state.
func (i *Interpreter) processMacro(tokens []objects.Token) uint8 {
	// TODO: Add import broker.
	tokens = tokens[1:]
	switch tokens[0].Type {
	case fract.TypeIf:
		return i.processMacroIf(tokens)
	case fract.TypeName:
		switch tokens[0].Value {
		case "define": // Macro variable.
			i.macroDefines = append(i.macroDefines, i.processMacroDefine(tokens))
		default:
			fract.Error(tokens[0], "Invalid macro!")
		}
	default:
		fract.Error(tokens[0], "Invalid macro!")
	}
	return fract.TypeNone
}
