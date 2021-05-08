package interpreter

import (
	"strconv"

	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// ProcessVariableSet Process variable set statement.
// tokens Tokens to process.
func (i *Interpreter) processVariableSet(tokens []objects.Token) {
	_name := tokens[0]

	// Name is not name?
	if _name.Type != fract.TypeName {
		fract.Error(_name, "This is not a valid name!")
	} else if _name.Value == grammar.TokenUnderscore {
		fract.Error(_name, "Ignore operator is cannot set!")
	}

	index, _ := i.varIndexByName(_name)
	if index == -1 {
		fract.Error(_name, "Variable is not defined in this name!: "+_name.Value)
	}

	variable := i.variables[index]

	// Check const state.
	if variable.Const {
		fract.Error(tokens[1], "Values is cannot changed of constant defines!")
	}

	setter := tokens[1]
	setIndex := -1

	// Array setter?
	if setter.Type == fract.TypeBrace && setter.Value == grammar.TokenLBracket {
		// Variable is not array?
		if !variable.Value.Array && variable.Value.Content[0].Type != fract.VALString {
			fract.Error(setter, "Variable is not array!")
		}

		// Find close bracket.
		for cindex := 2; cindex < len(tokens); cindex++ {
			current := tokens[cindex]
			if current.Type != fract.TypeBrace || current.Value != grammar.TokenRBracket {
				continue
			}

			valueList := vector.Sublist(tokens, 2, cindex-2)

			// Index value is empty?
			if valueList == nil {
				fract.Error(setter, "Index is not defined!")
			}

			position, err := strconv.Atoi(i.processValue(*valueList).Content[0].Data)
			if err != nil {
				fract.Error(setter, "Value out of range!")
			}

			if variable.Value.Array {
				position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
			} else {
				position = parser.ProcessArrayIndex(len(variable.Value.Content[0].Data), position)
			}
			if position == -1 {
				fract.Error(setter, "Index is out of range!")
			}

			setIndex = position
			vector.RemoveRange(&tokens, 1, cindex)
			setter = tokens[1]
			break
		}
	}

	/*
		// Setter is not a setter operator?
		if setter.Type != fract.TypeOperator && setter.Value != grammar.TokenEquals {
			fract.Error(setter, "This is not a setter operator!"+setter.Value)
		}*/

	// Value are not defined?
	if len(tokens) < 3 {
		fract.ErrorCustom(setter.File, setter.Line, setter.Column+len(setter.Value),
			"Value is not defined!")
	}

	value := i.processValue(*vector.Sublist(tokens, 2, len(tokens)-2))
	if value.Content == nil {
		fract.Error(tokens[2], "Invalid value!")
	}

	if setIndex != -1 {
		if value.Array {
			fract.Error(setter, "Array is cannot set as indexed value!")
		}

		switch setter.Value {
		case grammar.TokenEquals: // =
			if variable.Value.Array {
				variable.Value.Content[setIndex] = value.Content[0]
			} else {
				if value.Content[0].Type != fract.VALString {
					fract.Error(setter, "Value type is not string!")
				} else if len(value.Content[0].Data) > 1 {
					fract.Error(setter, "Value length is should be maximum one!")
				}
				bytes := []byte(variable.Value.Content[0].Data)
				if value.Content[0].Data == "" {
					bytes[setIndex] = 0
				} else {
					bytes[setIndex] = value.Content[0].Data[0]
				}
				variable.Value.Content[0].Data = string(bytes)
			}
		default: // Other assignments.
			if variable.Value.Array {
				variable.Value.Content[setIndex] = solveProcess(
					valueProcess{
						Operator: objects.Token{Value: string(setter.Value[:len(setter.Value)-1])},
						First:    tokens[0],
						FirstV: objects.Value{
							Content: []objects.DataFrame{variable.Value.Content[setIndex]},
						},
						Second:  setter,
						SecondV: value,
					}).Content[0]
			} else {
				value = solveProcess(
					valueProcess{
						Operator: objects.Token{Value: string(setter.Value[:len(setter.Value)-1])},
						First:    tokens[0],
						FirstV: objects.Value{
							Content: []objects.DataFrame{variable.Value.Content[setIndex]},
						},
						Second:  setter,
						SecondV: value,
					})
				if value.Content[0].Type != fract.VALString {
					fract.Error(setter, "Value type is not string!")
				} else if len(value.Content[0].Data) > 1 {
					fract.Error(setter, "Value length is should be maximum one!")
				}

				bytes := []byte(variable.Value.Content[0].Data)
				if value.Content[0].Data == "" {
					bytes[setIndex] = 0
				} else {
					bytes[setIndex] = value.Content[0].Data[0]
				}
				variable.Value.Content[0].Data = string(bytes)
			}
		}
	} else {
		switch setter.Value {
		case grammar.TokenEquals: // =
			variable.Value = value
		default: // Other assignments.
			variable.Value = solveProcess(
				valueProcess{
					Operator: objects.Token{Value: string(setter.Value[:len(setter.Value)-1])},
					First:    tokens[0],
					FirstV:   variable.Value,
					Second:   setter,
					SecondV:  value,
				})
		}
	}
}
