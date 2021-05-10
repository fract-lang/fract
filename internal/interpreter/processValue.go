package interpreter

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// valueProcess instance for solver.
type valueProcess struct {
	First    objects.Token // First value of process.
	FirstV   objects.Value // Value instance of first value.
	Second   objects.Token // Second value of process.
	SecondV  objects.Value // Value instance of second value.
	Operator objects.Token // Operator of process.
}

// processRange by value processor principles.
func (i *Interpreter) processRange(tokens *[]objects.Token) {
	for {
		_range, found := parser.DecomposeBrace(tokens, grammar.TokenLParenthes, grammar.TokenRParenthes, true)

		/* Parentheses are not found! */
		if found == -1 {
			return
		}

		val := i.processValue(_range)
		if val.Array {
			vector.Insert(tokens, found, objects.Token{
				Value: grammar.TokenLBracket,
				Type:  fract.TypeBrace,
			})
			for _, current := range val.Content {
				found++
				vector.Insert(tokens, found, objects.Token{
					Value: fract.FormatData(current),
					Type:  fract.TypeValue,
				})
				found++
				vector.Insert(tokens, found, objects.Token{
					Value: grammar.TokenComma,
					Type:  fract.TypeComma,
				})
			}
			found++
			vector.Insert(tokens, found, objects.Token{
				Value: grammar.TokenRBracket,
				Type:  fract.TypeBrace,
			})
		} else {
			if val.Content[0].Type == fract.VALString {
				vector.Insert(tokens, found, objects.Token{
					Value: grammar.TokenDoubleQuote + val.Content[0].Data + grammar.TokenDoubleQuote,
					Type:  fract.TypeValue,
				})
			} else {
				vector.Insert(tokens, found, objects.Token{
					Value: fract.FormatData(val.Content[0]),
					Type:  fract.TypeValue,
				})
			}
		}
	}
}

// solve process.
func solve(operator objects.Token, first, second float64) float64 {
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

	switch operator.Value {
	case grammar.TokenPlus: // Addition.
		result = first + second
	case grammar.TokenMinus: // Subtraction.
		result = first - second
	case grammar.TokenStar: // Multiply.
		result = first * second
	case
		grammar.TokenSlash,
		grammar.IntegerDivision: // Division.
		if first == 0 || second == 0 {
			fract.Error(operator, "Divide by zero!")
		}
		result = first / second

		if operator.Value == grammar.IntegerDivision {
			result = math.RoundToEven(result)
		}
	case grammar.TokenVerticalBar: // Binary or.
		result = float64(int(first) | int(second))
	case grammar.TokenAmper: // Binary and.
		result = float64(int(first) & int(second))
	case grammar.TokenCaret: // Bitwise exclusive or.
		result = float64(int(first) ^ int(second))
	case grammar.Exponentiation: // Exponentiation.
		result = math.Pow(first, second)
	case grammar.TokenPercent: // Mod.
		result = math.Mod(first, second)
	case grammar.LeftBinaryShift: // Left shift.
		if second < 0 {
			fract.Error(operator, "Shifter is cannot should be negative!")
		}
		result = float64(int(first) << int(second))
	case grammar.RightBinaryShift:
		if second < 0 {
			fract.Error(operator, "Shifter is cannot should be negative!")
		}
		result = float64(int(first) >> int(second))
	default:
		fract.Error(operator, "Operator is invalid!")
	}

	return result
}

// readyDataFrame to data.
func readyDataFrame(process valueProcess, dataFrame objects.DataFrame) objects.DataFrame {
	if process.FirstV.Content[0].Type == fract.VALString ||
		process.SecondV.Content[0].Type == fract.VALString {
		dataFrame.Type = fract.VALString
	} else if process.Operator.Value == grammar.TokenSlash ||
		process.Operator.Value == grammar.TokenBackslash ||
		process.FirstV.Content[0].Type == fract.VALFloat ||
		process.SecondV.Content[0].Type == fract.VALFloat {
		dataFrame.Type = fract.VALFloat
	}
	dataFrame.Data = fract.FormatData(dataFrame)
	return dataFrame
}

// solveProcess solve arithmetic process.
func solveProcess(process valueProcess) objects.Value {
	value := objects.Value{Content: []objects.DataFrame{{}}}

	// String?
	if (len(process.FirstV.Content) != 0 && process.FirstV.Content[0].Type == fract.VALString) ||
		(len(process.SecondV.Content) != 0 && process.SecondV.Content[0].Type == fract.VALString) {
		if process.FirstV.Content[0].Type == process.SecondV.Content[0].Type { // Both string?
			value.Content[0].Type = fract.VALString
			switch process.Operator.Value {
			case grammar.TokenPlus:
				value.Content[0].Data = process.FirstV.Content[0].Data + process.SecondV.Content[0].Data
			case grammar.TokenMinus:
				firstLen := len(process.FirstV.Content[0].Data)
				secondLen := len(process.SecondV.Content[0].Data)

				if firstLen == 0 || secondLen == 0 {
					value.Content[0].Data = ""
					break
				}

				if firstLen == 1 && secondLen > 1 {
					result, _ := strconv.ParseInt(process.FirstV.Content[0].Data, 10, 32)
					fRune := rune(result)
					for _, char := range process.SecondV.Content[0].Data {
						value.Content[0].Data += string(fRune - char)
					}
				} else if secondLen == 1 && firstLen > 1 {
					result, _ := strconv.ParseInt(process.SecondV.Content[0].Data, 10, 32)
					fRune := rune(result)
					for _, char := range process.FirstV.Content[0].Data {
						value.Content[0].Data += string(fRune - char)
					}
				} else {
					for index, char := range process.FirstV.Content[0].Data {
						value.Content[0].Data += string(char - rune(process.SecondV.Content[0].Data[index]))
					}
				}
			default:
				fract.Error(process.Operator, "This operator is not defined for string types!")
			}
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
					fract.Error(process.Second, "Array element count is not one or equals to first array!")
				}

				if strings.Contains(process.SecondV.Content[0].Data, grammar.TokenDot) {
					fract.Error(process.Second, "Only string and integer values cannot concatenate string values!")
				}

				result, _ := strconv.ParseInt(process.SecondV.Content[0].Data, 10, 64)
				_rune := rune(result)

				var sb strings.Builder
				for _, char := range process.FirstV.Content[0].Data {
					switch process.Operator.Value {
					case grammar.TokenPlus:
						sb.WriteByte(byte(char + _rune))
					case grammar.TokenMinus:
						sb.WriteByte(byte(char - _rune))
					default:
						fract.Error(process.Operator, "This operator is not defined for string types!")
					}
				}

				value.Content[0].Data = sb.String()
			} else {
				if process.SecondV.Content[0].Type != fract.VALInteger {
					fract.Error(process.Second, "Only string and integer values cannot concatenate string values!")
				}

				var sb strings.Builder
				result, _ := strconv.ParseInt(process.SecondV.Content[0].Data, 10, 64)
				_rune := rune(result)
				for _, char := range process.FirstV.Content[0].Data {
					switch process.Operator.Value {
					case grammar.TokenPlus:
						sb.WriteByte(byte(char + _rune))
					case grammar.TokenMinus:
						sb.WriteByte(byte(char - _rune))
					default:
						fract.Error(process.Operator, "This operator is not defined for string types!")
					}
				}

				value.Content[0].Data = sb.String()
			}
		} else {
			if process.FirstV.Array {
				if len(process.FirstV.Content) == 0 {
					value.Content = process.SecondV.Content
					return value
				}

				if len(process.FirstV.Content[0].Data) != len(process.SecondV.Content) &&
					(len(process.FirstV.Content[0].Data) != 1 && len(process.SecondV.Content) != 1) {
					fract.Error(process.Second, "Array element count is not one or equals to first array!")
				}

				if strings.Contains(process.FirstV.Content[0].Data, grammar.TokenDot) {
					fract.Error(process.Second, "Only string and integer values cannot concatenate string values!")
				}

				result, _ := strconv.ParseInt(process.FirstV.Content[0].Data, 10, 64)
				_rune := rune(result)

				var sb strings.Builder
				for _, char := range process.SecondV.Content[0].Data {
					switch process.Operator.Value {
					case grammar.TokenPlus:
						sb.WriteByte(byte(char + _rune))
					case grammar.TokenMinus:
						sb.WriteByte(byte(char - _rune))
					default:
						fract.Error(process.Operator, "This operator is not defined for string types!")
					}
				}

				value.Content[0].Data = sb.String()
			} else {
				if process.FirstV.Content[0].Type != fract.VALInteger {
					fract.Error(process.First, "Only string and integer values cannot concatenate string values!")
				}
				var sb strings.Builder
				result, _ := strconv.ParseInt(process.FirstV.Content[0].Data, 10, 64)
				_rune := rune(result)
				for _, char := range process.SecondV.Content[0].Data {
					switch process.Operator.Value {
					case grammar.TokenPlus:
						sb.WriteByte(byte(char + _rune))
					case grammar.TokenMinus:
						sb.WriteByte(byte(char - _rune))
					default:
						fract.Error(process.Operator, "This operator is not defined for string types!")
					}
				}

				value.Content[0].Data = sb.String()
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
			fract.Error(process.Second, "Array element count is not one or equals to first array!")
		}

		if len(process.FirstV.Content) == 1 {
			first := arithmetic.ToArithmetic(process.FirstV.Content[0].Data)
			for index, current := range process.SecondV.Content {
				process.SecondV.Content[index] = readyDataFrame(process,
					objects.DataFrame{
						Data: fmt.Sprintf(fract.FloatFormat,
							solve(process.Operator, first, arithmetic.ToArithmetic(current.Data))),
					})
			}
			value.Content = process.SecondV.Content
		} else if len(process.SecondV.Content) == 1 {
			second := arithmetic.ToArithmetic(process.SecondV.Content[0].Data)
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = readyDataFrame(process,
					objects.DataFrame{
						Data: fmt.Sprintf(fract.FloatFormat,
							solve(process.Operator, arithmetic.ToArithmetic(current.Data), second)),
					})
			}
			value.Content = process.FirstV.Content
		} else {
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = readyDataFrame(process,
					objects.DataFrame{
						Data: fmt.Sprintf(fract.FloatFormat,
							solve(process.Operator, arithmetic.ToArithmetic(current.Data),
								arithmetic.ToArithmetic(process.SecondV.Content[index].Data))),
					})
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
			process.FirstV.Content[index] = readyDataFrame(process,
				objects.DataFrame{
					Data: fmt.Sprintf(fract.FloatFormat,
						solve(process.Operator, arithmetic.ToArithmetic(current.Data), second)),
				})
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
			process.SecondV.Content[index] = readyDataFrame(process,
				objects.DataFrame{
					Data: fmt.Sprintf(fract.FloatFormat, solve(process.Operator, arithmetic.ToArithmetic(current.Data), first)),
				})
		}
		value.Content = process.SecondV.Content
	} else {
		if len(process.FirstV.Content) == 0 {
			process.FirstV.Content = []objects.DataFrame{{Data: "0"}}
		}

		value.Content[0] = readyDataFrame(process,
			objects.DataFrame{
				Data: fmt.Sprintf(fract.FloatFormat,
					solve(process.Operator, arithmetic.ToArithmetic(process.FirstV.Content[0].Data),
						arithmetic.ToArithmetic(process.SecondV.Content[0].Data))),
			})
	}

	return value
}

// applyMinus operator.
func applyMinus(minussed bool, value objects.Value) objects.Value {
	if !minussed {
		return value
	}

	val := objects.Value{
		Array:   value.Array,
		Content: append([]objects.DataFrame{}, value.Content...),
	}

	if val.Array {
		for index, data := range val.Content {
			if data.Type == fract.VALBoolean ||
				data.Type == fract.VALFloat ||
				data.Type == fract.VALInteger {
				data.Data = fmt.Sprintf(fract.FloatFormat, -arithmetic.ToArithmetic(data.Data))
				val.Content[index].Data = fract.FormatData(data)
			}
		}
		return val
	}

	if data := val.Content[0]; data.Type == fract.VALBoolean ||
		data.Type == fract.VALFloat ||
		data.Type == fract.VALInteger {
		data.Data = fmt.Sprintf(fract.FloatFormat, -arithmetic.ToArithmetic(data.Data))
		val.Content[0].Data = fract.FormatData(data)
	}

	return val
}

func (i *Interpreter) processOperationValue(first bool, operation *valueProcess, parts *[]objects.Token, index int) int {
	var (
		minussed bool
		token    = operation.First
	)

	if !first {
		token = operation.Second
	}

	minussed = token.Type == fract.TypeName && token.Value[0] == '-'

	if token.Type == fract.TypeName {
		if index < len(*parts)-1 {
			next := (*parts)[index+1]
			// Array?
			if next.Type == fract.TypeBrace {
				if next.Value == grammar.TokenLBracket {
					vindex, source := i.varIndexByName(token)
					if vindex == -1 {
						fract.Error(token, "Variable is not defined in this name!: "+token.Value)
					}

					// Find close bracket.
					cindex := index
					bracketCount := 0
					for ; cindex < len(*parts); cindex++ {
						current := (*parts)[cindex]
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

					valueList := vector.Sublist(*parts, index+2, cindex-index-2)
					// Index value is empty?
					if valueList == nil {
						fract.Error(token, "Index is not defined!")
					}

					value := i.processValue(*valueList)
					if value.Array {
						fract.Error((*parts)[index], "Arrays is not used in index access!")
					} else if value.Content[0].Type != fract.VALInteger {
						fract.Error((*parts)[index], "Only integer values can used in index access!")
					}
					position, err := strconv.Atoi(value.Content[0].Data)
					if err != nil {
						fract.Error((*parts)[index], "Invalid value!")
					}

					variable := source.variables[vindex]
					if !variable.Value.Array && variable.Value.Content[0].Type != fract.VALString {
						fract.Error((*parts)[index], "Index accessor is cannot used with non-array variables!")
					}

					if variable.Value.Array {
						position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
					} else {
						position = parser.ProcessArrayIndex(len(variable.Value.Content[0].Data), position)
					}

					if position == -1 {
						fract.Error((*parts)[index], "Index is out of range!")
					}

					vector.RemoveRange(parts, index+1, cindex-index)

					var data objects.DataFrame
					if variable.Value.Array {
						data = variable.Value.Content[position]
					} else {
						if variable.Value.Content[0].Type == fract.VALString {
							data = objects.DataFrame{
								Data: string(variable.Value.Content[0].Data[position]),
								Type: fract.VALString,
							}
						} else {
							data = objects.DataFrame{Data: fmt.Sprint(variable.Value.Content[0].Data[position])}
						}
					}

					if first {
						operation.FirstV.Content = []objects.DataFrame{data}
						operation.FirstV.Array = false
						operation.FirstV = applyMinus(minussed, operation.FirstV)
					} else {
						operation.SecondV.Content = []objects.DataFrame{data}
						operation.SecondV.Array = false
						operation.SecondV = applyMinus(minussed, operation.SecondV)
					}

					return 0
				} else if next.Value == grammar.TokenLParenthes { // Function?
					// Find close parentheses.
					cindex := index + 1
					bracketCount := 1
					for ; cindex < len(*parts); cindex++ {
						current := (*parts)[cindex]
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
					value := i.processFunctionCall(*vector.Sublist(*parts, index, cindex-index))
					if !operation.FirstV.Array && value.Content == nil {
						fract.Error(token, "Function is not return any value!")
					}
					vector.RemoveRange(parts, index+1, cindex-index-1)
					value = applyMinus(minussed, value)
					if first {
						operation.FirstV = value
					} else {
						operation.SecondV = value
					}
					return 0
				}
			}
		}

		vindex, source := i.varIndexByName(token)
		if vindex == -1 {
			fract.Error(token, "Variable is not defined in this name!: "+token.Value)
		}

		variable := source.variables[vindex]

		if first {
			operation.FirstV = applyMinus(minussed, variable.Value)
		} else {
			operation.SecondV = applyMinus(minussed, variable.Value)
		}
		return 0
	} else if token.Type == fract.TypeBrace {
		if token.Value == grammar.TokenRBracket {
			// Find open bracket.
			bracketCount := 1
			oindex := index - 1
			for ; oindex >= 0; oindex-- {
				current := (*parts)[oindex]
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
			if oindex == 0 || (*parts)[oindex-1].Type != fract.TypeName {
				if first {
					operation.FirstV.Array = true
					operation.FirstV.Content = i.processArrayValue(*vector.Sublist(*parts, oindex, index-oindex+1)).Content
					operation.FirstV = applyMinus(minussed, operation.FirstV)
				} else {
					operation.SecondV.Array = true
					operation.SecondV.Content = i.processArrayValue(*vector.Sublist(*parts, oindex, index-oindex+1)).Content
					operation.SecondV = applyMinus(minussed, operation.SecondV)
				}
				vector.RemoveRange(parts, oindex, index-oindex)
				return index - oindex
			}

			endToken := (*parts)[oindex-1]
			vindex, source := i.varIndexByName(endToken)
			if vindex == -1 {
				fract.Error(endToken, "Variable is not defined in this name!: "+endToken.Value)
			}
			valueList := vector.Sublist(*parts, oindex+1, index-oindex-1)
			// Index value is empty?
			if valueList == nil {
				fract.Error(endToken, "Index is not defined!")
			}

			value := i.processValue(*valueList)
			if value.Array {
				fract.Error((*parts)[index], "Arrays is not used in index access!")
			} else if value.Content[0].Type != fract.VALInteger {
				fract.Error((*parts)[index], "Only integer values can used in index access!")
			}

			position, err := strconv.Atoi(value.Content[0].Data)
			if err != nil {
				fract.Error((*parts)[oindex], "Invalid value!")
			}

			variable := source.variables[vindex]

			if !variable.Value.Array && variable.Value.Content[0].Type != fract.VALString {
				fract.Error((*parts)[oindex], "Index accessor is cannot used with non-array variables!")
			}

			if variable.Value.Array {
				position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
			} else {
				position = parser.ProcessArrayIndex(len(variable.Value.Content[0].Data), position)
			}

			if position == -1 {
				fract.Error((*parts)[oindex], "Index is out of range!")
			}
			vector.RemoveRange(parts, oindex-1, index-oindex+1)

			var data objects.DataFrame
			if variable.Value.Array {
				data = variable.Value.Content[position]
			} else {
				if variable.Value.Content[0].Type == fract.VALString {
					data = objects.DataFrame{
						Data: string(variable.Value.Content[0].Data[position]),
						Type: fract.VALString,
					}
				} else {
					data = objects.DataFrame{Data: fmt.Sprint(variable.Value.Content[0].Data[position])}
				}
			}

			if first {
				operation.FirstV.Content = []objects.DataFrame{data}
				operation.FirstV.Array = false
				operation.FirstV = applyMinus(minussed, operation.FirstV)
			} else {
				operation.SecondV.Content = []objects.DataFrame{data}
				operation.FirstV.Array = false
				operation.SecondV = applyMinus(minussed, operation.SecondV)
			}

			return index - oindex + 1
		} else if token.Value == grammar.TokenLBracket {
			// Array initializer.

			// Find close brace.
			cindex := index + 1
			braceCount := 1
			for ; cindex < len(*parts); cindex++ {
				current := (*parts)[cindex]
				if current.Type == fract.TypeBrace {
					if current.Value == grammar.TokenLBracket {
						braceCount++
					} else if current.Value == grammar.TokenRBracket {
						braceCount--
						if braceCount == 0 {
							break
						}
					}
				}
			}

			value := applyMinus(minussed, i.processArrayValue(*vector.Sublist(*parts, index, cindex-index+1)))
			if first {
				operation.FirstV = value
			} else {
				operation.SecondV = value
			}
			vector.RemoveRange(parts, index+1, cindex-index)
			return 0
		} else if token.Value == grammar.TokenRParenthes {
			// Function.

			// Find open parentheses.
			bracketCount := 1
			oindex := index - 1
			for ; oindex >= 0; oindex-- {
				current := (*parts)[oindex]
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
			value := i.processFunctionCall(*vector.Sublist(*parts, oindex, index-oindex+1))
			if value.Content == nil {
				fract.Error((*parts)[oindex], "Function is not return any value!")
			}
			value = applyMinus(minussed, value)
			if first {
				operation.FirstV = value
			} else {
				operation.SecondV = value
			}
			vector.RemoveRange(parts, oindex, index-oindex)
			return index - oindex
		}
	}

	//
	// Single value.
	//

	if (token.Type == fract.TypeValue &&
		token.Value != grammar.KwTrue &&
		token.Value != grammar.KwFalse) &&
		!strings.HasPrefix(token.Value, grammar.TokenQuote) &&
		!strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) {
		if strings.Contains(token.Value, grammar.TokenDot) ||
			strings.ContainsAny(token.Value, "eE") {
			token.Type = fract.VALFloat
		} else {
			token.Type = fract.VALInteger
		}

		if token.Value != "NaN" {
			prs, _ := new(big.Float).SetString(token.Value)
			val, _ := prs.Float64()
			token.Value = fmt.Sprint(val)
		}
	}

	if first {
		operation.FirstV.Array = false
		if strings.HasPrefix(token.Value, grammar.TokenQuote) ||
			strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) { // String?
			operation.FirstV.Content = []objects.DataFrame{{
				Data: token.Value[1 : len(token.Value)-1],
				Type: fract.VALString,
			}}
			token.Type = fract.TypeNone // Skip type check.
		} else {
			operation.FirstV.Content = []objects.DataFrame{{Data: token.Value}}
		}
	} else {
		operation.SecondV.Array = false
		if strings.HasPrefix(token.Value, grammar.TokenQuote) ||
			strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) { // String?
			operation.SecondV.Content = []objects.DataFrame{{
				Data: token.Value[1 : len(token.Value)-1],
				Type: fract.VALString,
			}}
			token.Type = fract.TypeNone // Skip type check.
		} else {
			operation.SecondV.Content = []objects.DataFrame{{Data: token.Value}}
		}
	}

	// Type check.
	if token.Type != fract.TypeNone {
		if token.Type == fract.TypeBooleanTrue ||
			token.Type == fract.TypeBooleanFalse {
			if first {
				operation.FirstV.Content[0].Type = fract.VALBoolean
				operation.FirstV = applyMinus(minussed, operation.FirstV)
			} else {
				operation.SecondV.Content[0].Type = fract.VALBoolean
				operation.SecondV = applyMinus(minussed, operation.SecondV)
			}
		} else if token.Type == fract.VALFloat { // Float?
			if first {
				operation.FirstV.Content[0].Type = fract.VALFloat
				operation.FirstV = applyMinus(minussed, operation.FirstV)
			} else {
				operation.SecondV.Content[0].Type = fract.VALFloat
				operation.SecondV = applyMinus(minussed, operation.SecondV)
			}
		}
	}

	return 0
}

func (i *Interpreter) processArrayValue(tokens []objects.Token) objects.Value {
	value := objects.Value{
		Content: []objects.DataFrame{},
		Array:   true,
	}

	first := tokens[0]

	comma := 1
	brace := 0
	for index := 1; index < len(tokens)-1; index++ {
		if current := tokens[index]; current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket ||
				current.Value == grammar.TokenLParenthes {
				brace++
			} else {
				brace--
			}
		} else if current.Type == fract.TypeComma && brace == 0 {
			lst := vector.Sublist(tokens, comma, index-comma)
			if lst == nil {
				fract.Error(first, "Value is not defined!")
			}
			val := i.processValue(*lst)
			value.Content = append(value.Content, val.Content...)
			comma = index + 1
		}
	}

	if comma < len(tokens)-1 {
		lst := vector.Sublist(tokens, comma, len(tokens)-comma-1)
		if lst == nil {
			fract.Error(first, "Value is not defined!")
		}
		val := i.processValue(*lst)
		value.Content = append(value.Content, val.Content...)
	}

	return value
}

func (i *Interpreter) processValue(tokens []objects.Token) objects.Value {
	value := objects.Value{
		Content: []objects.DataFrame{{}},
		Array:   false,
	}

	i.processRange(&tokens)

	// Is conditional expression?
	brace := 0
	for _, current := range tokens {
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket ||
				current.Value == grammar.TokenLParenthes {
				brace++
			} else {
				brace--
			}
		} else if brace == 0 && current.Type == fract.TypeOperator &&
			(current.Value == grammar.LogicalAnd || current.Value == grammar.LogicalOr ||
				current.Value == grammar.Equals || current.Value == grammar.NotEquals ||
				current.Value == grammar.TokenGreat || current.Value == grammar.TokenLess ||
				current.Value == grammar.GreaterEquals || current.Value == grammar.LessEquals) {
			value.Content = []objects.DataFrame{{
				Data: i.processCondition(tokens),
				Type: fract.VALBoolean,
			}}
			return value
		}
	}

	parts := parser.DecomposeArithmeticProcesses(tokens)

	if priorityIndex := parser.IndexProcessPriority(*parts); priorityIndex != -1 {
		// Decompose arithmetic operations.
		for priorityIndex != -1 {
			var operation valueProcess
			operation.First = (*parts)[priorityIndex-1]
			priorityIndex -= i.processOperationValue(true, &operation, parts, priorityIndex-1)

			operation.Operator = (*parts)[priorityIndex]

			operation.Second = (*parts)[priorityIndex+1]
			priorityIndex -= i.processOperationValue(false, &operation, parts, priorityIndex+1)

			resultValue := solveProcess(operation)

			operation.Operator.Value = grammar.TokenPlus
			operation.Second = (*parts)[priorityIndex+1]
			operation.FirstV = value
			operation.SecondV = resultValue

			resultValue = solveProcess(operation)
			value = resultValue

			// Remove processed processes.
			vector.RemoveRange(parts, priorityIndex-1, 3)
			vector.Insert(parts, priorityIndex-1, objects.Token{Value: "0"})

			// Find next operator.
			priorityIndex = parser.IndexProcessPriority(*parts)
		}
	} else {
		var operation valueProcess
		operation.First = (*parts)[0]
		operation.FirstV.Array = true //* Ignore nil control if function call.
		i.processOperationValue(true, &operation, parts, 0)
		value = operation.FirstV
	}

	return value
}
