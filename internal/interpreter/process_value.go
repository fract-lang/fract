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

// Get string arithmetic compatible data.
func arith(token objects.Token, d objects.Data) string {
	ret := d.String()
	switch d.Type {
	case objects.VALFunction:
		fract.Error(token, "\""+ret+"\" is not compatible with arithmetic processes!")
	}
	return ret
}

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
		_range, found := parser.DecomposeBrace(tokens, "(", ")", true)
		/* Parentheses are not found! */
		if found == -1 {
			return
		}
		val := i.processValue(_range)
		if val.Array {
			vector.Insert(tokens, found, objects.Token{
				Value: "[",
				Type:  fract.TypeBrace,
			})
			for _, current := range val.Content {
				found++
				vector.Insert(tokens, found, objects.Token{
					Value: current.Format(),
					Type:  fract.TypeValue,
				})
				found++
				vector.Insert(tokens, found, objects.Token{
					Value: ",",
					Type:  fract.TypeComma,
				})
			}
			found++
			vector.Insert(tokens, found, objects.Token{
				Value: "]",
				Type:  fract.TypeBrace,
			})
		} else {
			if val.Content[0].Type == objects.VALString {
				vector.Insert(tokens, found, objects.Token{
					Value: "\"" + val.Content[0].String() + "\"",
					Type:  fract.TypeValue,
				})
			} else {
				vector.Insert(tokens, found, objects.Token{
					Value: val.Content[0].Format(),
					Type:  fract.TypeValue,
					//! Add another fields for panic.
					Line:   _range[0].Line,
					Column: _range[0].Column,
					File:   _range[0].File,
				})
			}
		}
	}
}

// solve process.
func solve(operator objects.Token, first, second float64) float64 {
	var result float64
	if operator.Value == "\\" ||
		operator.Value == grammar.IntegerDivideWithBigger { // Divide with bigger.
		if operator.Value == "\\" {
			operator.Value = "/"
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
	case "+": // Addition.
		result = first + second
	case "-": // Subtraction.
		result = first - second
	case "*": // Multiply.
		result = first * second
	case "/", grammar.IntegerDivision: // Division.
		if first == 0 || second == 0 {
			fract.Error(operator, "Divide by zero!")
		}
		result = first / second
		if operator.Value == grammar.IntegerDivision {
			result = math.RoundToEven(result)
		}
	case "|": // Binary or.
		result = float64(int(first) | int(second))
	case "&": // Binary and.
		result = float64(int(first) & int(second))
	case "^": // Bitwise exclusive or.
		result = float64(int(first) ^ int(second))
	case grammar.Exponentiation: // Exponentiation.
		result = math.Pow(first, second)
	case "%": // Mod.
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

// readyData to data.
func readyData(process valueProcess, data objects.Data) objects.Data {
	if process.FirstV.Content[0].Type == objects.VALString || process.SecondV.Content[0].Type == objects.VALString {
		data.Type = objects.VALString
	} else if process.Operator.Value == "/" || process.Operator.Value == "\\" ||
		process.FirstV.Content[0].Type == objects.VALFloat || process.SecondV.Content[0].Type == objects.VALFloat {
		data.Data = data.Format()
		return data
	}
	return data
}

// solveProcess solve arithmetic process.
func solveProcess(process valueProcess) objects.Value {
	value := objects.Value{Content: []objects.Data{{}}}
	// String?
	if (len(process.FirstV.Content) != 0 && process.FirstV.Content[0].Type == objects.VALString) ||
		(len(process.SecondV.Content) != 0 && process.SecondV.Content[0].Type == objects.VALString) {
		if process.FirstV.Content[0].Type == process.SecondV.Content[0].Type { // Both string?
			value.Content[0].Type = objects.VALString
			switch process.Operator.Value {
			case "+":
				value.Content[0].Data = process.FirstV.Content[0].String() + process.SecondV.Content[0].String()
			case "-":
				firstLen := len(process.FirstV.Content[0].String())
				secondLen := len(process.SecondV.Content[0].String())
				if firstLen == 0 || secondLen == 0 {
					value.Content[0].Data = ""
					break
				}
				if firstLen == 1 && secondLen > 1 {
					result, _ := strconv.ParseInt(process.FirstV.Content[0].String(), 10, 32)
					fRune := rune(result)
					for _, char := range process.SecondV.Content[0].String() {
						value.Content[0].Data = value.Content[0].String() + string(fRune-char)
					}
				} else if secondLen == 1 && firstLen > 1 {
					result, _ := strconv.ParseInt(process.SecondV.Content[0].String(), 10, 32)
					fRune := rune(result)
					for _, char := range process.FirstV.Content[0].String() {
						value.Content[0].Data = value.Content[0].String() + string(fRune-char)
					}
				} else {
					for index, char := range process.FirstV.Content[0].String() {
						value.Content[0].Data = value.Content[0].String() + string(char-rune(process.SecondV.Content[0].String()[index]))
					}
				}
			default:
				fract.Error(process.Operator, "This operator is not defined for string types!")
			}
			return value
		}

		value.Content[0].Type = objects.VALString
		if process.FirstV.Content[0].Type == objects.VALString {
			if process.SecondV.Array {
				if len(process.SecondV.Content) == 0 {
					value.Content = process.FirstV.Content
					return value
				}
				if len(process.FirstV.Content[0].String()) != len(process.SecondV.Content) &&
					(len(process.FirstV.Content[0].String()) != 1 && len(process.SecondV.Content) != 1) {
					fract.Error(process.Second, "Array element count is not one or equals to first array!")
				}
				if strings.Contains(process.SecondV.Content[0].String(), ".") {
					fract.Error(process.Second, "Only string and integer values cannot concatenate string values!")
				}
				result, _ := strconv.ParseInt(process.SecondV.Content[0].String(), 10, 64)
				_rune := rune(result)
				var sb strings.Builder
				for _, char := range process.FirstV.Content[0].String() {
					switch process.Operator.Value {
					case "+":
						sb.WriteByte(byte(char + _rune))
					case "-":
						sb.WriteByte(byte(char - _rune))
					default:
						fract.Error(process.Operator, "This operator is not defined for string types!")
					}
				}
				value.Content[0].Data = sb.String()
			} else {
				if process.SecondV.Content[0].Type != objects.VALInteger {
					fract.Error(process.Second, "Only string and integer values cannot concatenate string values!")
				}
				var sb strings.Builder
				result, _ := strconv.ParseInt(process.SecondV.Content[0].String(), 10, 64)
				_rune := rune(result)
				for _, char := range process.FirstV.Content[0].String() {
					switch process.Operator.Value {
					case "+":
						sb.WriteByte(byte(char + _rune))
					case "-":
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
				if len(process.FirstV.Content[0].String()) != len(process.SecondV.Content) &&
					(len(process.FirstV.Content[0].String()) != 1 && len(process.SecondV.Content) != 1) {
					fract.Error(process.Second, "Array element count is not one or equals to first array!")
				}
				if strings.Contains(process.FirstV.Content[0].String(), ".") {
					fract.Error(process.Second, "Only string and integer values cannot concatenate string values!")
				}
				result, _ := strconv.ParseInt(process.FirstV.Content[0].String(), 10, 64)
				_rune := rune(result)
				var sb strings.Builder
				for _, char := range process.SecondV.Content[0].String() {
					switch process.Operator.Value {
					case "+":
						sb.WriteByte(byte(char + _rune))
					case "-":
						sb.WriteByte(byte(char - _rune))
					default:
						fract.Error(process.Operator, "This operator is not defined for string types!")
					}
				}
				value.Content[0].Data = sb.String()
			} else {
				if process.FirstV.Content[0].Type != objects.VALInteger {
					fract.Error(process.First, "Only string and integer values cannot concatenate string values!")
				}
				var sb strings.Builder
				result, _ := strconv.ParseInt(process.FirstV.Content[0].String(), 10, 64)
				_rune := rune(result)
				for _, char := range process.SecondV.Content[0].String() {
					switch process.Operator.Value {
					case "+":
						sb.WriteByte(byte(char + _rune))
					case "-":
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
			first := arithmetic.ToArithmetic(arith(process.Operator, process.FirstV.Content[0]))
			for index, current := range process.SecondV.Content {
				process.SecondV.Content[index] = readyData(process,
					objects.Data{
						Data: fmt.Sprintf(fract.FloatFormat, solve(process.Operator, first, arithmetic.ToArithmetic(arith(process.Operator, current)))),
					})
			}
			value.Content = process.SecondV.Content
		} else if len(process.SecondV.Content) == 1 {
			second := arithmetic.ToArithmetic(arith(process.Operator, process.SecondV.Content[0]))
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = readyData(process,
					objects.Data{
						Data: fmt.Sprintf(fract.FloatFormat, solve(process.Operator, arithmetic.ToArithmetic(arith(process.Operator, current)), second)),
					})
			}
			value.Content = process.FirstV.Content
		} else {
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = readyData(process,
					objects.Data{
						Data: fmt.Sprintf(fract.FloatFormat, solve(process.Operator, arithmetic.ToArithmetic(arith(process.Operator, current)), arithmetic.ToArithmetic(process.SecondV.Content[index].String()))),
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

		second := arithmetic.ToArithmetic(arith(process.Operator, process.SecondV.Content[0]))
		for index, current := range process.FirstV.Content {
			process.FirstV.Content[index] = readyData(process,
				objects.Data{
					Data: fmt.Sprintf(fract.FloatFormat, solve(process.Operator, arithmetic.ToArithmetic(arith(process.Operator, current)), second)),
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

		first := arithmetic.ToArithmetic(arith(process.Operator, process.FirstV.Content[0]))
		for index, current := range process.SecondV.Content {
			process.SecondV.Content[index] = readyData(process,
				objects.Data{
					Data: fmt.Sprintf(fract.FloatFormat, solve(process.Operator, arithmetic.ToArithmetic(arith(process.Operator, current)), first)),
				})
		}
		value.Content = process.SecondV.Content
	} else {
		if len(process.FirstV.Content) == 0 {
			process.FirstV.Content = []objects.Data{{Data: "0"}}
		}
		value.Content[0] = readyData(process,
			objects.Data{
				Data: fmt.Sprintf(fract.FloatFormat,
					solve(process.Operator,
						arithmetic.ToArithmetic(arith(process.Operator, process.FirstV.Content[0])),
						arithmetic.ToArithmetic(arith(process.Operator, process.SecondV.Content[0])))),
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
		Content: append([]objects.Data{}, value.Content...),
	}
	if val.Array {
		for index, data := range val.Content {
			if data.Type == objects.VALBoolean ||
				data.Type == objects.VALFloat ||
				data.Type == objects.VALInteger {
				data.Data = fmt.Sprintf(fract.FloatFormat, -arithmetic.ToArithmetic(data.String()))
				val.Content[index].Data = data.Format()
			}
		}
		return val
	}
	if data := val.Content[0]; data.Type == objects.VALBoolean ||
		data.Type == objects.VALFloat ||
		data.Type == objects.VALInteger {
		data.Data = fmt.Sprintf(fract.FloatFormat, -arithmetic.ToArithmetic(data.String()))
		val.Content[0].Data = data.Format()
	}
	return val
}

func (i *Interpreter) processOperationValue(first bool, operation *valueProcess, parts *[]objects.Token, index int) int {
	var (
		token  = operation.First
		result = &operation.FirstV
	)

	if !first {
		token = operation.Second
		result = &operation.SecondV
	}

	minussed := token.Type == fract.TypeName && token.Value[0] == '-'
	if token.Type == fract.TypeName {
		if index < len(*parts)-1 {
			next := (*parts)[index+1]
			// Array?
			if next.Type == fract.TypeBrace {
				if next.Value == "[" {
					vindex, t, source := i.DefineByName(token) //! HERE
					if vindex == -1 || t != 'v' {
						fract.Error(token, "Variable is not defined in this name: "+token.Value)
					}
					// Find close bracket.
					cindex := index
					bracketCount := 0
					for ; cindex < len(*parts); cindex++ {
						current := (*parts)[cindex]
						if current.Type == fract.TypeBrace {
							if current.Value == "[" {
								bracketCount++
							} else if current.Value == "]" {
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
					} else if value.Content[0].Type != objects.VALInteger {
						fract.Error((*parts)[index], "Only integer values can used in index access!")
					}
					position, err := strconv.Atoi(arith((*valueList)[0], value.Content[0]))
					if err != nil {
						fract.Error((*parts)[index], "Invalid value!")
					}
					variable := source.variables[vindex]
					if !variable.Value.Array && variable.Value.Content[0].Type != objects.VALString {
						fract.Error((*parts)[index], "Index accessor is cannot used with non-array variables!")
					}
					if variable.Value.Array {
						position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
					} else {
						position = parser.ProcessArrayIndex(len(variable.Value.Content[0].String()), position)
					}
					if position == -1 {
						fract.Error((*parts)[index], "Index is out of range!")
					}
					vector.RemoveRange(parts, index+1, cindex-index)
					var data objects.Data
					if variable.Value.Array {
						data = variable.Value.Content[position]
					} else {
						if variable.Value.Content[0].Type == objects.VALString {
							data = objects.Data{
								Data: string(variable.Value.Content[0].String()[position]),
								Type: objects.VALString,
							}
						} else {
							data = objects.Data{Data: fmt.Sprint(variable.Value.Content[0].String()[position])}
						}
					}
					result.Content = []objects.Data{data}
					result.Array = false
					*result = applyMinus(minussed, *result)
					return 0
				} else if next.Value == "(" { // Function?
					// Find close parentheses.
					cindex := index + 1
					bracketCount := 1
					for ; cindex < len(*parts); cindex++ {
						current := (*parts)[cindex]
						if current.Type == fract.TypeBrace {
							if current.Value == "(" {
								bracketCount++
							} else if current.Value == ")" {
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
					*result = applyMinus(minussed, value)
					return 0
				}
			}
		}
		vindex, t, source := i.DefineByName(token) //! HERE
		if vindex == -1 {
			fract.Error(token, "Variable is not defined in this name: "+token.Value)
		}
		switch t {
		case 'f':
			*result = objects.Value{
				Content: []objects.Data{{
					Data: source.functions[vindex],
					Type: objects.VALFunction,
				}},
			}
		case 'v':
			variable := source.variables[vindex]
			value := variable.Value
			if !variable.Mutable { //! Immutability.
				value.Content = append(make([]objects.Data, 0), variable.Value.Content...)
			}

			*result = applyMinus(minussed, value)
		}
		return 0
	} else if token.Type == fract.TypeBrace {
		if token.Value == "}" {
			// Find open bracket.
			bracketCount := 1
			oindex := index - 1
			for ; oindex >= 0; oindex-- {
				current := (*parts)[oindex]
				if current.Type == fract.TypeBrace {
					if current.Value == "}" {
						bracketCount++
					} else if current.Value == "{" {
						bracketCount--
						if bracketCount == 0 {
							break
						}
					}
				}
			}
			// Finished?
			if oindex == 0 || (*parts)[oindex-1].Type != fract.TypeName {
				result.Array = true
				result.Content = i.processArrayValue(*vector.Sublist(*parts, oindex, index-oindex+1)).Content
				*result = applyMinus(minussed, *result)
				vector.RemoveRange(parts, oindex, index-oindex)
				return index - oindex
			}
			endToken := (*parts)[oindex-1]
			vindex, t, source := i.DefineByName(token) //! HERE
			if vindex == -1 || t != 'v' {
				fract.Error(endToken, "Variable is not defined in this name: "+endToken.Value)
			}
			valueList := vector.Sublist(*parts, oindex+1, index-oindex-1)
			// Index value is empty?
			if valueList == nil {
				fract.Error(endToken, "Index is not defined!")
			}
			value := i.processValue(*valueList)
			if value.Array {
				fract.Error((*parts)[index], "Arrays is not used in index access!")
			} else if value.Content[0].Type != objects.VALInteger {
				fract.Error((*parts)[index], "Only integer values can used in index access!")
			}
			position, err := strconv.Atoi(arith((*valueList)[0], value.Content[0]))
			if err != nil {
				fract.Error((*parts)[oindex], "Invalid value!")
			}
			variable := source.variables[vindex]
			if !variable.Value.Array && variable.Value.Content[0].Type != objects.VALString {
				fract.Error((*parts)[oindex], "Index accessor is cannot used with non-array variables!")
			}
			if variable.Value.Array {
				position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
			} else {
				position = parser.ProcessArrayIndex(len(variable.Value.Content[0].String()), position)
			}
			if position == -1 {
				fract.Error((*parts)[oindex], "Index is out of range!")
			}
			vector.RemoveRange(parts, oindex-1, index-oindex+1)
			var data objects.Data
			if variable.Value.Array {
				data = variable.Value.Content[position]
			} else {
				if variable.Value.Content[0].Type == objects.VALString {
					data = objects.Data{
						Data: string(variable.Value.Content[0].String()[position]),
						Type: objects.VALString,
					}
				} else {
					data = objects.Data{Data: fmt.Sprint(variable.Value.Content[0].String()[position])}
				}
			}
			result.Content = []objects.Data{data}
			result.Array = false
			*result = applyMinus(minussed, *result)
			return index - oindex + 1
		} else if token.Value == "[" {
			// Array initializer.

			// Find close brace.
			cindex := index + 1
			braceCount := 1
			for ; cindex < len(*parts); cindex++ {
				current := (*parts)[cindex]
				if current.Type == fract.TypeBrace {
					if current.Value == "[" {
						braceCount++
					} else if current.Value == "]" {
						braceCount--
						if braceCount == 0 {
							break
						}
					}
				}
			}
			*result = applyMinus(minussed, i.processArrayValue(*vector.Sublist(*parts, index, cindex-index+1)))
			vector.RemoveRange(parts, index+1, cindex-index)
			return 0
		} else if token.Value == "]" {
			// Function.

			// Find open parentheses.
			bracketCount := 1
			oindex := index - 1
			for ; oindex >= 0; oindex-- {
				current := (*parts)[oindex]
				if current.Type == fract.TypeBrace {
					if current.Value == "]" {
						bracketCount++
					} else if current.Value == "[" {
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
			*result = applyMinus(minussed, value)
			vector.RemoveRange(parts, oindex, index-oindex)
			return index - oindex
		}
	}

	//* Single value.
	if strings.HasPrefix(token.Value, "object.") {
		fract.Error(token, "\""+token.Value+"\" is not compatible with arithmetic processes!")
	}
	if (token.Type == fract.TypeValue && token.Value != grammar.KwTrue && token.Value != grammar.KwFalse) &&
		token.Value[0] != '\'' && token.Value[0] != '"' {
		if strings.Contains(token.Value, ".") || strings.ContainsAny(token.Value, "eE") {
			token.Type = objects.VALFloat
		} else {
			token.Type = objects.VALInteger
		}
		if token.Value != grammar.KwNaN {
			prs, _ := new(big.Float).SetString(token.Value)
			val, _ := prs.Float64()
			token.Value = fmt.Sprint(val)
		}
	}
	result.Array = false
	if token.Value[0] == '\'' || token.Value[0] == '"' { // String?
		result.Content = []objects.Data{{
			Data: token.Value[1 : len(token.Value)-1],
			Type: objects.VALString,
		}}
		token.Type = fract.TypeNone // Skip type check.
	} else {
		result.Content = []objects.Data{{Data: token.Value}}
	}
	//* Type check.
	if token.Type != fract.TypeNone {
		if token.Value == grammar.KwTrue || token.Value == grammar.KwFalse {
			result.Content[0].Type = objects.VALBoolean
			*result = applyMinus(minussed, *result)
		} else if token.Type == objects.VALFloat { // Float?
			result.Content[0].Type = objects.VALFloat
			*result = applyMinus(minussed, *result)
		}
	}
	return 0
}

func (i *Interpreter) processArrayValue(tokens []objects.Token) objects.Value {
	value := objects.Value{
		Array:   true,
		Content: []objects.Data{},
	}
	first := tokens[0]
	comma := 1
	brace := 0
	for index := 1; index < len(tokens)-1; index++ {
		if current := tokens[index]; current.Type == fract.TypeBrace {
			if current.Value == "[" || current.Value == "{" || current.Value == "(" {
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
	i.processRange(&tokens)
	value := objects.Value{Content: []objects.Data{{}}}

	// Is conditional expression?
	brace := 0
	for _, current := range tokens {
		if current.Type == fract.TypeBrace {
			if current.Value == "{" || current.Value == "[" || current.Value == "(" {
				brace++
			} else {
				brace--
			}
		} else if brace == 0 && current.Type == fract.TypeOperator &&
			(current.Value == grammar.LogicalAnd || current.Value == grammar.LogicalOr || current.Value == grammar.Equals ||
				current.Value == grammar.NotEquals || current.Value == ">" || current.Value == "<" ||
				current.Value == grammar.GreaterEquals || current.Value == grammar.LessEquals) {
			value.Content = []objects.Data{{
				Data: i.processCondition(tokens),
				Type: objects.VALBoolean,
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
			operation.Operator.Value = "+"
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
