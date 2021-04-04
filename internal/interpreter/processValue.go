/*
	processValue Function
*/

package interpreter

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// valueProcess Value process instance.
type valueProcess struct {
	// First value of process.
	First obj.Token
	// Value instance of first value.
	FirstV obj.Value
	// Second value of process.
	Second obj.Token
	// Value instance of second value.
	SecondV obj.Value
	// Operator of process.
	Operator obj.Token
}

// processRange Process range by value processor principles.
// tokens Tokens to process.
func (i *Interpreter) processRange(tokens *[]obj.Token) {
	/* Check parentheses range. */
	for {
		_range, found := parser.DecomposeBrace(tokens, grammar.TokenLParenthes,
			grammar.TokenRParenthes, true)

		/* Parentheses are not found! */
		if found == -1 {
			return
		}

		val := i.processValue(&_range)
		if val.Array {
			vector.Insert(tokens, found, obj.Token{
				Value: grammar.TokenLBracket,
				Type:  fract.TypeBrace,
			})
			for _, current := range val.Content {
				found++
				vector.Insert(tokens, found, obj.Token{
					Value: current.Data,
					Type:  fract.TypeValue,
				})
				found++
				vector.Insert(tokens, found, obj.Token{
					Value: grammar.TokenComma,
					Type:  fract.TypeComma,
				})
			}
			found++
			vector.Insert(tokens, found, obj.Token{
				Value: grammar.TokenRBracket,
				Type:  fract.TypeBrace,
			})
		} else {
			if val.Content[0].Type == fract.VALString {
				vector.Insert(tokens, found, obj.Token{
					Value: grammar.TokenDoubleQuote + val.Content[0].Data + grammar.TokenDoubleQuote,
					Type:  fract.TypeValue,
				})
				continue
			}
			vector.Insert(tokens, found, obj.Token{
				Value: val.Content[0].Data,
				Type:  fract.TypeValue,
			})
		}
	}
}

// solve process.
// operator Operator of process.
// first First value.
// second Second value.
func solve(operator obj.Token, first, second float64) float64 {
	var result float64

	if operator.Value == grammar.TokenBackslash ||
		operator.Value == grammar.IntegerDivideWithBigger { // Divide with bigger.
		if operator.Value == grammar.TokenBackslash {
			operator.Value = grammar.TokenSlash
		} else {
			operator.Value = grammar.IntegerDivision
		}

		if first < second {
			cache := first
			first = second
			second = cache
		}
	}

	if operator.Value == grammar.TokenPlus { // Addition.
		result = first + second
	} else if operator.Value == grammar.TokenMinus { // Subtraction.
		result = first - second
	} else if operator.Value == grammar.TokenStar { // Multiply.
		result = first * second
	} else if operator.Value == grammar.TokenSlash ||
		operator.Value == grammar.IntegerDivision { // Division.
		if first == 0 || second == 0 {
			fract.Error(operator, "Divide by zero!")
		}
		result = first / second

		if operator.Value == grammar.IntegerDivision {
			result = math.RoundToEven(result)
		}
	} else if operator.Value == grammar.TokenCaret { // Exponentiation.
		result = math.Pow(first, second)
	} else if operator.Value == grammar.TokenPercent { // Mod.
		result = math.Mod(first, second)
	} else if operator.Value == grammar.LeftShift { // Left shift.
		if second < 0 {
			fract.Error(operator, "Shifter is can not should be negative!")
		}
		result = float64(int64(first) << int64(second))
	} else if operator.Value == grammar.RightShift { // Right shift.
		if second < 0 {
			fract.Error(operator, "Shifter is can not should be negative!")
		}
		result = float64(int64(first) >> int64(second))
	} else {
		fract.Error(operator, "Operator is invalid!")
	}

	return result
}

// solveProcess Solve arithmetic process.
// process Process to solve.
func solveProcess(process valueProcess) obj.Value {
	value := obj.Value{Content: []obj.DataFrame{{}}}

	// String?
	if (len(process.FirstV.Content) != 0 && process.FirstV.Content[0].Type == fract.VALString) ||
		(len(process.SecondV.Content) != 0 && process.SecondV.Content[0].Type == fract.VALString) {
		if len(process.FirstV.Content) == 0 {
			value.Content = process.SecondV.Content
			return value
		} else if len(process.SecondV.Content) == 0 {
			value.Content = process.FirstV.Content
			return value
		}

		if process.FirstV.Content[0].Type == process.SecondV.Content[0].Type { // Both string?
			value.Content = []obj.DataFrame{{
				Data: process.FirstV.Content[0].Data + process.SecondV.Content[0].Data,
				Type: fract.VALString,
			}}
			return value
		}

		value.Content[0].Type = fract.VALString

		if process.FirstV.Content[0].Type == fract.VALString {
			if process.SecondV.Array {
				if len(process.SecondV.Content) == 0 {
					value.Content = process.FirstV.Content
					return value
				}
				if len(process.FirstV.Content[0].Data) != len(process.SecondV.Content) &&
					(len(process.FirstV.Content[0].Data) != 1 && len(process.SecondV.Content) != 1) {
					fract.Error(process.Second,
						"Array element count is not one or equals to first array!")
				}
				var sb strings.Builder
				for _, char := range process.FirstV.Content[0].Data {
					if strings.Contains(process.SecondV.Content[0].Data, grammar.TokenDot) {
						fract.Error(process.Second, "Float values cannot concatenate string values!")
					}
					sb.WriteRune(char + rune(arithmetic.ToArithmetic(process.SecondV.Content[0].Data)))
				}
				value.Content = []obj.DataFrame{
					{
						Data: sb.String(),
						Type: fract.VALString,
					},
				}
			} else {
				if process.SecondV.Content[0].Type == fract.VALFloat {
					fract.Error(process.Second, "Float values cannot concatenate string values!")
				}
				var sb strings.Builder
				val := rune(arithmetic.ToArithmetic(process.SecondV.Content[0].Data))
				for _, char := range process.FirstV.Content[0].Data {
					sb.WriteRune(char + val)
				}
				value.Content = []obj.DataFrame{
					{
						Data: sb.String(),
						Type: fract.VALString,
					},
				}
			}
		} else {
			if process.FirstV.Array {
				if len(process.FirstV.Content) == 0 {
					value.Content = process.SecondV.Content
					return value
				}
				if len(process.FirstV.Content[0].Data) != len(process.SecondV.Content) &&
					(len(process.FirstV.Content[0].Data) != 1 && len(process.SecondV.Content) != 1) {
					fract.Error(process.Second,
						"Array element count is not one or equals to first array!")
				}
				var sb strings.Builder
				for _, char := range process.SecondV.Content[0].Data {
					if strings.Contains(process.FirstV.Content[0].Data, grammar.TokenDot) {
						fract.Error(process.Second, "Float values cannot concatenate string values!")
					}
					sb.WriteRune(char + rune(arithmetic.ToArithmetic(process.FirstV.Content[0].Data)))
				}
				value.Content = []obj.DataFrame{
					{
						Data: sb.String(),
						Type: fract.VALString,
					},
				}
			} else {
				if process.FirstV.Content[0].Type == fract.VALFloat {
					fract.Error(process.First, "Float values cannot concatenate string values!")
				}
				var sb strings.Builder
				val := rune(arithmetic.ToArithmetic(process.FirstV.Content[0].Data))
				for _, char := range process.SecondV.Content[0].Data {
					sb.WriteRune(char + val)
				}
				value.Content = []obj.DataFrame{
					{
						Data: sb.String(),
						Type: fract.VALString,
					},
				}
			}
		}
		return value
	}

	// ****************************

	if process.FirstV.Array && process.SecondV.Array {
		value.Array = true

		if len(process.FirstV.Content) == 0 {
			value.Content = process.SecondV.Content
			return value
		} else if len(process.SecondV.Content) == 0 {
			value.Content = process.FirstV.Content
			return value
		}

		if len(process.FirstV.Content) != len(process.SecondV.Content) &&
			(len(process.FirstV.Content) != 1 && len(process.SecondV.Content) != 1) {
			fract.Error(process.Second,
				"Array element count is not one or equals to first array!")
		}

		if len(process.FirstV.Content) == 1 {
			first := arithmetic.ToArithmetic(process.FirstV.Content[0].Data)
			for index, current := range process.SecondV.Content {
				process.SecondV.Content[index] = obj.DataFrame{
					Data: fmt.Sprintf("%g",
						solve(process.Operator, first, arithmetic.ToArithmetic(current.Data)))}
			}
			value.Content = process.SecondV.Content
		} else if len(process.SecondV.Content) == 1 {
			second := arithmetic.ToArithmetic(process.SecondV.Content[0].Data)
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = obj.DataFrame{
					Data: fmt.Sprintf("%g",
						solve(process.Operator, arithmetic.ToArithmetic(current.Data), second))}
			}
			value.Content = process.FirstV.Content
		} else {
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = obj.DataFrame{
					Data: fmt.Sprintf("%g",
						solve(process.Operator, arithmetic.ToArithmetic(current.Data),
							arithmetic.ToArithmetic(process.SecondV.Content[index].Data)))}
			}
			value.Content = process.FirstV.Content
		}
	} else if process.FirstV.Array {
		value.Array = true

		if len(process.FirstV.Content) == 0 {
			value.Content = process.SecondV.Content
			return value
		} else if len(process.SecondV.Content) == 0 {
			value.Content = process.FirstV.Content
			return value
		}

		second := arithmetic.ToArithmetic(process.SecondV.Content[0].Data)
		for index, current := range process.FirstV.Content {
			process.FirstV.Content[index] = obj.DataFrame{
				Data: fmt.Sprintf("%g",
					solve(process.Operator, arithmetic.ToArithmetic(current.Data), second))}
		}
		value.Content = process.FirstV.Content
	} else if process.SecondV.Array {
		value.Array = true

		if len(process.FirstV.Content) == 0 {
			value.Content = process.SecondV.Content
			return value
		} else if len(process.SecondV.Content) == 0 {
			value.Content = process.FirstV.Content
			return value
		}

		first := arithmetic.ToArithmetic(process.FirstV.Content[0].Data)
		for index, current := range process.SecondV.Content {
			process.SecondV.Content[index] = obj.DataFrame{
				Data: fmt.Sprintf("%g",
					solve(process.Operator, arithmetic.ToArithmetic(current.Data), first))}
		}
		value.Content = process.SecondV.Content
	} else {
		if len(process.FirstV.Content) == 0 {
			process.FirstV.Content = []obj.DataFrame{{Data: "0"}}
		}

		value.Content = []obj.DataFrame{{
			Data: fmt.Sprintf("%g",
				solve(process.Operator, arithmetic.ToArithmetic(process.FirstV.Content[0].Data),
					arithmetic.ToArithmetic(process.SecondV.Content[0].Data)))}}
	}

	return value
}

// __processValue Process value.
// first This is first value.
// token Token to process.
// tokens Tokens to process.
// index Index of token.
func (i *Interpreter) _processValue(first bool, operation *valueProcess,
	tokens *[]obj.Token, index int) int {
	token := operation.First
	if !first {
		token = operation.Second
	}

	if token.Type == fract.TypeName {
		if index < len(*tokens)-1 {
			next := (*tokens)[index+1]
			// Array?
			if next.Type == fract.TypeBrace {
				if next.Value == grammar.TokenLBracket {
					vindex := i.varIndexByName(token.Value)
					if vindex == -1 {
						fract.Error(token, "Variable is not defined in this name!: "+token.Value)
					}

					// Find close bracket.
					cindex := index + 1
					bracketCount := 1
					for ; cindex < len(*tokens); cindex++ {
						current := (*tokens)[cindex]
						if current.Type == fract.TypeBrace {
							if current.Value == grammar.TokenLBracket {
								bracketCount++
							} else if current.Value == grammar.TokenRBracket {
								bracketCount--
								if bracketCount == 0 {
									break
								}
							}
						}
					}

					valueList := vector.Sublist(*tokens, index+2, cindex-index-3)
					// Index value is empty?
					if valueList == nil {
						fract.Error(token, "Index is not defined!")
					}

					value := i.processValue(valueList)
					if value.Array {
						fract.Error((*tokens)[index], "Arrays is not used in index access!")
					} else if value.Content[0].Type != fract.VALInteger {
						fract.Error((*tokens)[index],
							"Only integer values can used in index access!")
					}
					position, err := strconv.Atoi(value.Content[0].Data)
					if err != nil {
						fract.Error((*tokens)[index], "Value out of range!")
					}

					variable := i.variables[vindex]

					if !variable.Value.Array && variable.Value.Content[0].Type != fract.VALString {
						fract.Error((*tokens)[index],
							"Index accessor is cannot used with non-array variables!")
					}

					if variable.Value.Array {
						position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
					} else {
						position = parser.ProcessArrayIndex(len(variable.Value.Content[0].Data), position)
					}

					if position == -1 {
						fract.Error((*tokens)[index], "Index is out of range!")
					}
					vector.RemoveRange(tokens, index+1, cindex-index-1)

					var data obj.DataFrame
					if variable.Value.Array {
						data = variable.Value.Content[position]
					} else {
						data = obj.DataFrame{
							Data: fmt.Sprintf("%d", variable.Value.Content[0].Data[position]),
							Type: fract.VALInteger,
						}
					}

					if first {
						operation.FirstV.Content = []obj.DataFrame{data}
						operation.FirstV.Array = false
					} else {
						operation.SecondV.Content = []obj.DataFrame{data}
						operation.SecondV.Array = false
					}
					return 0
				} else if next.Value == grammar.TokenLParenthes { // Function?
					// Find close parentheses.
					cindex := index + 1
					bracketCount := 1
					for ; cindex < len(*tokens); cindex++ {
						current := (*tokens)[cindex]
						if current.Type == fract.TypeBrace {
							if current.Value == grammar.TokenLParenthes {
								bracketCount++
							} else if current.Value == grammar.TokenRParenthes {
								bracketCount--
								if bracketCount == 0 {
									break
								}
							}
						}
					}
					value := i.processFunctionCall(*vector.Sublist(*tokens, index, cindex-index))
					if !operation.FirstV.Array && value.Content == nil {
						fract.Error(token, "Function is not return any value!")
					}
					vector.RemoveRange(tokens, index+1, cindex-index-1)
					if first {
						operation.FirstV = value
					} else {
						operation.SecondV = value
					}
					return 0
				}
			}
		}

		vindex := i.varIndexByName(token.Value)
		if vindex == -1 {
			fract.Error(token, "Variable is not defined in this name!: "+token.Value)
		}

		variable := i.variables[vindex]

		if first {
			operation.FirstV = variable.Value
		} else {
			operation.SecondV = variable.Value
		}
		return 0
	} else if token.Type == fract.TypeBrace {
		if token.Value == grammar.TokenRBracket {
			// Find open bracket.
			bracketCount := 1
			oindex := index - 1
			for ; oindex >= 0; oindex-- {
				current := (*tokens)[oindex]
				if current.Type == fract.TypeBrace {
					if current.Value == grammar.TokenRBracket {
						bracketCount++
					} else if current.Value == grammar.TokenLBracket {
						bracketCount--
						if bracketCount == 0 {
							break
						}
					}
				}
			}

			// Finished?
			if oindex == 0 {
				if first {
					operation.FirstV.Array = true
					operation.FirstV.Content = i.processArrayValue(
						vector.Sublist(*tokens, oindex, index-oindex+1)).Content
				} else {
					operation.SecondV.Array = true
					operation.SecondV.Content = i.processArrayValue(
						vector.Sublist(*tokens, oindex, index-oindex+1)).Content
				}
				vector.RemoveRange(tokens, oindex, index-oindex)
				return index - oindex
			}

			endToken := (*tokens)[oindex-1]
			vindex := i.varIndexByName(endToken.Value)
			if vindex == -1 {
				fract.Error(endToken, "Variable is not defined in this name!: "+endToken.Value)
			}
			valueList := vector.Sublist(*tokens, oindex+1, index-oindex-1)
			// Index value is empty?
			if valueList == nil {
				fract.Error(endToken, "Index is not defined!")
			}

			value := i.processValue(valueList)
			if value.Array {
				fract.Error((*tokens)[index],
					"Arrays is not used in index access!")
			} else if value.Content[0].Type != fract.VALInteger {
				fract.Error((*tokens)[index],
					"Only integer values can used in index access!")
			}

			position, err := strconv.Atoi(value.Content[0].Data)
			if err != nil {
				fract.Error((*tokens)[oindex], "Value out of range!")
			}

			variable := i.variables[vindex]

			if !variable.Value.Array && variable.Value.Content[0].Type != fract.VALString {
				fract.Error((*tokens)[oindex],
					"Index accessor is cannot used with non-array variables!")
			}

			if variable.Value.Array {
				position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
			} else {
				position = parser.ProcessArrayIndex(len(variable.Value.Content[0].Data), position)
			}

			if position == -1 {
				fract.Error((*tokens)[oindex], "Index is out of range!")
			}
			vector.RemoveRange(tokens, oindex-1, index-oindex+1)

			var data obj.DataFrame
			if variable.Value.Array {
				data = variable.Value.Content[position]
			} else {
				data = obj.DataFrame{
					Data: fmt.Sprintf("%d", variable.Value.Content[0].Data[position]),
					Type: fract.VALInteger,
				}
			}

			if first {
				operation.FirstV.Content = []obj.DataFrame{data}
				operation.FirstV.Array = false
				if variable.Value.Content[0].Type == fract.VALString {
					operation.FirstV.Content[0].Type = fract.VALString
				}
			} else {
				operation.SecondV.Content = []obj.DataFrame{data}
				operation.FirstV.Array = false
			}

			return index - oindex + 1
		} else if token.Value == grammar.TokenLBracket {
			// Array initializer.

			// Find close brace.
			cindex := index + 1
			braceCount := 1
			for ; cindex < len(*tokens); cindex++ {
				current := (*tokens)[cindex]
				if current.Type == fract.TypeBrace {
					if current.Value == grammar.TokenLBracket {
						fract.Error(current, "Arrays is cannot take array value as element!")
					} else if current.Value == grammar.TokenRBracket {
						braceCount--
						if braceCount == 0 {
							break
						}
					}
				}
			}

			value := i.processArrayValue(vector.Sublist(*tokens, index, cindex-index+1))
			if first {
				operation.FirstV = value
			} else {
				operation.SecondV = value
			}
			vector.RemoveRange(tokens, index+1, cindex-index-1)
			return 0
		} else if token.Value == grammar.TokenRParenthes {
			// Function.

			// Find open parentheses.
			bracketCount := 1
			oindex := index - 1
			for ; oindex >= 0; oindex-- {
				current := (*tokens)[oindex]
				if current.Type == fract.TypeBrace {
					if current.Value == grammar.TokenRBracket {
						bracketCount++
					} else if current.Value == grammar.TokenLBracket {
						bracketCount--
						if bracketCount == 0 {
							break
						}
					}
				}
			}
			oindex++
			value := i.processFunctionCall(*vector.Sublist(*tokens, oindex, index-oindex+1))
			if value.Content == nil {
				fract.Error((*tokens)[oindex], "Function is not return any value!")
			}
			if first {
				operation.FirstV = value
			} else {
				operation.SecondV = value
			}
			vector.RemoveRange(tokens, oindex, index-oindex)
			return index - oindex
		}
	}

	//
	// Single value.
	//

	if token.Type != fract.TypeBooleanTrue &&
		token.Type != fract.TypeBooleanFalse &&
		!strings.HasPrefix(token.Value, grammar.TokenQuote) &&
		!strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) {
		val, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			fract.Error(token, "Value out of range!")
		}
		token.Value = fmt.Sprintf("%g", val)
	}

	if first {
		operation.FirstV.Array = false
		if strings.HasPrefix(token.Value, grammar.TokenQuote) ||
			strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) { // String?
			operation.FirstV.Content = []obj.DataFrame{{
				Data: token.Value[1 : len(token.Value)-1],
				Type: fract.VALString,
			}}
		} else {
			operation.FirstV.Content = []obj.DataFrame{{Data: token.Value}}
		}
	} else {
		operation.SecondV.Array = false
		if strings.HasPrefix(token.Value, grammar.TokenQuote) ||
			strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) { // String?
			operation.SecondV.Content = []obj.DataFrame{{
				Data: token.Value[1 : len(token.Value)-1],
				Type: fract.VALString,
			}}
		} else {
			operation.SecondV.Content = []obj.DataFrame{{Data: token.Value}}
		}
	}

	// Type check.
	if token.Type == fract.TypeBooleanTrue ||
		token.Type == fract.TypeBooleanFalse {
		if first {
			operation.FirstV.Content[0].Type = fract.VALBoolean
		} else {
			operation.SecondV.Content[0].Type = fract.VALBoolean
		}
	} else if strings.Contains(token.Value, grammar.TokenDot) { // Float?
		if first {
			operation.FirstV.Content[0].Type = fract.VALFloat
		} else {
			operation.SecondV.Content[0].Type = fract.VALFloat
		}
	}

	return 0
}

// processArrayValue Process array value.
// tokens Tokens.
func (i *Interpreter) processArrayValue(tokens *[]obj.Token) obj.Value {
	value := obj.Value{
		Content: []obj.DataFrame{},
		Array:   true,
	}

	first := (*tokens)[0]

	comma := 1
	for index := 1; index < len(*tokens)-1; index++ {
		current := (*tokens)[index]
		if current.Type == fract.TypeComma {
			lst := vector.Sublist(*tokens, comma, index-comma)
			if lst == nil {
				fract.Error(first, "Value is not defined!")
			}
			val := i.processValue(lst)
			value.Content = append(value.Content, val.Content...)
			comma = index + 1
		}
	}

	if comma < len(*tokens)-1 {
		lst := vector.Sublist(*tokens, comma, len(*tokens)-comma-1)
		if lst == nil {
			fract.Error(first, "Value is not defined!")
		}
		val := i.processValue(lst)
		value.Content = append(value.Content, val.Content...)
	}

	return value
}

// processValue Process value.
// tokens Tokens to process.
func (i *Interpreter) processValue(tokens *[]obj.Token) obj.Value {
	value := obj.Value{
		Content: []obj.DataFrame{{}},
		Array:   false,
	}

	i.processRange(tokens)

	// Is conditional expression?
	brace := 0
	for _, current := range *tokens {
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket ||
				current.Value == grammar.TokenLParenthes {
				brace++
			} else {
				brace--
			}
		} else if brace == 0 && current.Type == fract.TypeOperator &&
			(current.Value == grammar.TokenAmper || current.Value == grammar.TokenVerticalBar ||
				current.Value == grammar.Equals || current.Value == grammar.NotEquals ||
				current.Value == grammar.TokenGreat || current.Value == grammar.TokenLess ||
				current.Value == grammar.GreaterEquals || current.Value == grammar.LessEquals) {
			value.Content = []obj.DataFrame{{
				Data: i.processCondition(tokens),
				Type: fract.VALBoolean,
			}}
			return value
		}
	}

	if priorityIndex := parser.IndexProcessPriority(*tokens); priorityIndex != -1 {
		// Decompose arithmetic operations.
		for priorityIndex != -1 {
			var operation valueProcess
			operation.First = (*tokens)[priorityIndex-1]
			priorityIndex -= i._processValue(true, &operation,
				tokens, priorityIndex-1)
			operation.Operator = (*tokens)[priorityIndex]

			operation.Second = (*tokens)[priorityIndex+1]
			priorityIndex -= i._processValue(false, &operation,
				tokens, priorityIndex+1)

			resultValue := solveProcess(operation)

			operation.Operator.Value = grammar.TokenPlus
			operation.Second = (*tokens)[priorityIndex+1]
			operation.FirstV = value
			operation.SecondV = resultValue

			resultValue = solveProcess(operation)
			value = resultValue

			// Remove processed processes.
			vector.RemoveRange(tokens, priorityIndex-1, 3)
			vector.Insert(tokens, priorityIndex-1, obj.Token{Value: "0"})

			// Find next operator.
			priorityIndex = parser.IndexProcessPriority(*tokens)
		}
	} else {
		var operation valueProcess
		operation.First = (*tokens)[0]
		operation.FirstV.Array = true // Ignore nil control if function call.
		i._processValue(true, &operation, tokens, 0)
		value = operation.FirstV
	}

	return value
}
