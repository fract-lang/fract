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
		_range, found := parser.DecomposeBrace(tokens, grammar.TokenLParenthes, grammar.TokenRParenthes, true)

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
					Value: fract.FormatData(current),
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
			} else {
				vector.Insert(tokens, found, obj.Token{
					Value: fract.FormatData(val.Content[0]),
					Type:  fract.TypeValue,
				})
			}
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
		result = float64(int64(first) | int64(second))
	case grammar.TokenAmper: // Binary and.
		result = float64(int64(first) & int64(second))
	case grammar.TokenCaret: // Bitwise exclusive or.
		result = float64(int64(first) ^ int64(second))
	case grammar.Exponentiation: // Exponentiation.
		result = math.Pow(first, second)
	case grammar.TokenPercent: // Mod.
		result = math.Mod(first, second)
	case grammar.LeftBinaryShift: // Left shift.
		if second < 0 {
			fract.Error(operator, "Shifter is cannot should be negative!")
		}
		result = float64(int64(first) << int64(second))
	case grammar.RightBinaryShift:
		if second < 0 {
			fract.Error(operator, "Shifter is cannot should be negative!")
		}
		result = float64(int64(first) >> int64(second))
	default:
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

				result, _ := strconv.ParseInt(process.SecondV.Content[0].Data, 10, 32)
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
				result, _ := strconv.ParseInt(process.SecondV.Content[0].Data, 10, 32)
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

				result, _ := strconv.ParseInt(process.FirstV.Content[0].Data, 10, 32)
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
				result, _ := strconv.ParseInt(process.FirstV.Content[0].Data, 10, 32)
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

	// readyDataFrame DataFrame ready to data.
	// dataFrame Destination dataframe.
	readyDataFrame := func(dataFrame obj.DataFrame) obj.DataFrame {
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
				process.SecondV.Content[index] = readyDataFrame(obj.DataFrame{
					Data: fmt.Sprintf(fract.FloatFormat,
						solve(process.Operator, first, arithmetic.ToArithmetic(current.Data)))})
			}
			value.Content = process.SecondV.Content
		} else if len(process.SecondV.Content) == 1 {
			second := arithmetic.ToArithmetic(process.SecondV.Content[0].Data)
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = readyDataFrame(obj.DataFrame{
					Data: fmt.Sprintf(fract.FloatFormat,
						solve(process.Operator, arithmetic.ToArithmetic(current.Data), second))})
			}
			value.Content = process.FirstV.Content
		} else {
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = readyDataFrame(obj.DataFrame{
					Data: fmt.Sprintf(fract.FloatFormat,
						solve(process.Operator, arithmetic.ToArithmetic(current.Data),
							arithmetic.ToArithmetic(process.SecondV.Content[index].Data)))})
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
			process.FirstV.Content[index] = readyDataFrame(obj.DataFrame{
				Data: fmt.Sprintf(fract.FloatFormat,
					solve(process.Operator, arithmetic.ToArithmetic(current.Data), second))})
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
			process.SecondV.Content[index] = readyDataFrame(obj.DataFrame{
				Data: fmt.Sprintf(fract.FloatFormat,
					solve(process.Operator, arithmetic.ToArithmetic(current.Data), first))})
		}
		value.Content = process.SecondV.Content
	} else {
		if len(process.FirstV.Content) == 0 {
			process.FirstV.Content = []obj.DataFrame{{Data: "0"}}
		}

		value.Content[0] = readyDataFrame(obj.DataFrame{
			Data: fmt.Sprintf(fract.FloatFormat,
				solve(process.Operator, arithmetic.ToArithmetic(process.FirstV.Content[0].Data),
					arithmetic.ToArithmetic(process.SecondV.Content[0].Data)))})
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
	var (
		minussed bool
		token    = operation.First
	)

	// applyMinus Apply minus assignment.
	// value Value to apply.
	applyMinus := func(value obj.Value) obj.Value {
		if !minussed {
			return value
		}

		val := obj.Value{
			Array:   value.Array,
			Content: append([]obj.DataFrame{}, value.Content...),
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

	if !first {
		token = operation.Second
	}

	minussed = token.Type == fract.TypeName && token.Value[0] == '-'

	if token.Type == fract.TypeName {
		if index < len(*tokens)-1 {
			next := (*tokens)[index+1]
			// Array?
			if next.Type == fract.TypeBrace {
				if next.Value == grammar.TokenLBracket {
					vindex, source := i.varIndexByName(token)
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
						fract.Error((*tokens)[index], "Only integer values can used in index access!")
					}
					position, err := strconv.Atoi(value.Content[0].Data)
					if err != nil {
						fract.Error((*tokens)[index], "Invalid value!")
					}

					variable := source.variables[vindex]
					if !variable.Value.Array && variable.Value.Content[0].Type != fract.VALString {
						fract.Error((*tokens)[index], "Index accessor is cannot used with non-array variables!")
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
						if variable.Value.Content[0].Type == fract.VALString {
							data = obj.DataFrame{
								Data: string(variable.Value.Content[0].Data[position]),
								Type: fract.VALString,
							}
						} else {
							data = obj.DataFrame{Data: fmt.Sprint(variable.Value.Content[0].Data[position])}
						}
					}

					if first {
						operation.FirstV.Content = []obj.DataFrame{data}
						operation.FirstV.Array = false
						operation.FirstV = applyMinus(operation.FirstV)
					} else {
						operation.SecondV.Content = []obj.DataFrame{data}
						operation.SecondV.Array = false
						operation.SecondV = applyMinus(operation.SecondV)
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
					value = applyMinus(value)
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
			operation.FirstV = applyMinus(variable.Value)
		} else {
			operation.SecondV = applyMinus(variable.Value)
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
					operation.FirstV = applyMinus(operation.FirstV)
				} else {
					operation.SecondV.Array = true
					operation.SecondV.Content = i.processArrayValue(
						vector.Sublist(*tokens, oindex, index-oindex+1)).Content
					operation.SecondV = applyMinus(operation.SecondV)
				}
				vector.RemoveRange(tokens, oindex, index-oindex)
				return index - oindex
			}

			endToken := (*tokens)[oindex-1]
			vindex, source := i.varIndexByName(endToken)
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
				fract.Error((*tokens)[index], "Arrays is not used in index access!")
			} else if value.Content[0].Type != fract.VALInteger {
				fract.Error((*tokens)[index], "Only integer values can used in index access!")
			}

			position, err := strconv.Atoi(value.Content[0].Data)
			if err != nil {
				fract.Error((*tokens)[oindex], "Invalid value!")
			}

			variable := source.variables[vindex]

			if !variable.Value.Array && variable.Value.Content[0].Type != fract.VALString {
				fract.Error((*tokens)[oindex], "Index accessor is cannot used with non-array variables!")
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
				if variable.Value.Content[0].Type == fract.VALString {
					data = obj.DataFrame{
						Data: string(variable.Value.Content[0].Data[position]),
						Type: fract.VALString,
					}
				} else {
					data = obj.DataFrame{Data: fmt.Sprint(variable.Value.Content[0].Data[position])}
				}
			}

			if first {
				operation.FirstV.Content = []obj.DataFrame{data}
				operation.FirstV.Array = false
				operation.FirstV = applyMinus(operation.FirstV)
			} else {
				operation.SecondV.Content = []obj.DataFrame{data}
				operation.FirstV.Array = false
				operation.SecondV = applyMinus(operation.SecondV)
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
						braceCount++
					} else if current.Value == grammar.TokenRBracket {
						braceCount--
						if braceCount == 0 {
							break
						}
					}
				}
			}

			value := i.processArrayValue(vector.Sublist(*tokens, index, cindex-index+1))
			value = applyMinus(value)
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
			value = applyMinus(value)
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

	if (token.Type == fract.TypeValue &&
		token.Value != grammar.KwTrue &&
		token.Value != grammar.KwFalse) &&
		!strings.HasPrefix(token.Value, grammar.TokenQuote) &&
		!strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) {
		if strings.Contains(token.Value, grammar.TokenDot) ||
			strings.Contains(token.Value, "e") {
			val, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				fract.Error(token, "Invalid value!")
			}
			token.Value = fmt.Sprintf(fract.FloatFormat, val)
		} else {
			val, err := strconv.ParseInt(token.Value, 10, 64)
			if err != nil {
				fract.Error(token, "Invalid value!")
			}
			token.Value = fmt.Sprint(val)
		}
	}

	if first {
		operation.FirstV.Array = false
		if strings.HasPrefix(token.Value, grammar.TokenQuote) ||
			strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) { // String?
			operation.FirstV.Content = []obj.DataFrame{{
				Data: token.Value[1 : len(token.Value)-1],
				Type: fract.VALString,
			}}
			token.Type = fract.TypeNone // Skip type check.
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
			token.Type = fract.TypeNone // Skip type check.
		} else {
			operation.SecondV.Content = []obj.DataFrame{{Data: token.Value}}
		}
	}

	// Type check.
	if token.Type != fract.TypeNone {
		if token.Type == fract.TypeBooleanTrue ||
			token.Type == fract.TypeBooleanFalse {
			if first {
				operation.FirstV.Content[0].Type = fract.VALBoolean
				operation.FirstV = applyMinus(operation.FirstV)
			} else {
				operation.SecondV.Content[0].Type = fract.VALBoolean
				operation.SecondV = applyMinus(operation.SecondV)
			}
		} else if strings.Contains(token.Value, grammar.TokenDot) { // Float?
			if first {
				operation.FirstV.Content[0].Type = fract.VALFloat
				operation.FirstV = applyMinus(operation.FirstV)
			} else {
				operation.SecondV.Content[0].Type = fract.VALFloat
				operation.SecondV = applyMinus(operation.SecondV)
			}
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
	brace := 0
	for index := 1; index < len(*tokens)-1; index++ {
		if current := (*tokens)[index]; current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket ||
				current.Value == grammar.TokenLParenthes {
				brace++
			} else {
				brace--
			}
		} else if current.Type == fract.TypeComma && brace == 0 {
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
			(current.Value == grammar.LogicalAnd || current.Value == grammar.LogicalOr ||
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

	parts := parser.DecomposeArithmeticProcesses(*tokens)

	if priorityIndex := parser.IndexProcessPriority(*parts); priorityIndex != -1 {
		// Decompose arithmetic operations.
		for priorityIndex != -1 {
			var operation valueProcess
			operation.First = (*parts)[priorityIndex-1]
			priorityIndex -= i._processValue(true, &operation, parts, priorityIndex-1)
			operation.Operator = (*parts)[priorityIndex]

			operation.Second = (*parts)[priorityIndex+1]
			priorityIndex -= i._processValue(false, &operation, parts, priorityIndex+1)

			resultValue := solveProcess(operation)

			operation.Operator.Value = grammar.TokenPlus
			operation.Second = (*parts)[priorityIndex+1]
			operation.FirstV = value
			operation.SecondV = resultValue

			resultValue = solveProcess(operation)
			value = resultValue

			// Remove processed processes.
			vector.RemoveRange(parts, priorityIndex-1, 3)
			vector.Insert(parts, priorityIndex-1, obj.Token{Value: "0"})

			// Find next operator.
			priorityIndex = parser.IndexProcessPriority(*parts)
		}
	} else {
		var operation valueProcess
		operation.First = (*parts)[0]
		operation.FirstV.Array = true //* Ignore nil control if function call.
		i._processValue(true, &operation, parts, 0)
		value = operation.FirstV
	}

	return value
}
