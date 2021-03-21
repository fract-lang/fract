/*
	processValue Function
*/

package interpreter

import (
	"fmt"
	"math"
	"strings"

	"github.com/fract-lang/fract/pkg/arithmetic"
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/parser"
	"github.com/fract-lang/fract/pkg/vector"
)

// valueProcess Value process instance.
type valueProcess struct {
	// First value of process.
	First objects.Token
	// Value instance of first value.
	FirstV objects.Value
	// Second value of process.
	Second objects.Token
	// Value instance of second value.
	SecondV objects.Value
	// Operator of process.
	Operator objects.Token
}

// processRange Process range by value processor principles.
// tokens Tokens to process.
func (i *Interpreter) processRange(tokens *vector.Vector) {
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
			tokens.Insert(found, objects.Token{
				Value: grammar.TokenLBrace,
				Type:  fract.TypeBrace,
			})
			for _, current := range val.Content {
				found++
				tokens.Insert(found, objects.Token{
					Value: current,
					Type:  fract.TypeValue,
				})
				found++
				tokens.Insert(found, objects.Token{
					Value: grammar.TokenComma,
					Type:  fract.TypeComma,
				})
			}
			found++
			tokens.Insert(found, objects.Token{
				Value: grammar.TokenRBrace,
				Type:  fract.TypeBrace,
			})
		} else {
			if val.Type == fract.VALString {
				tokens.Insert(found, objects.Token{
					Value: grammar.TokenDoubleQuote + val.Content[0] + grammar.TokenDoubleQuote,
					Type:  fract.TypeValue,
				})
				continue
			}
			tokens.Insert(found, objects.Token{
				Value: val.Content[0],
				Type:  fract.TypeValue,
			})
		}
	}
}

// solveProcess Solve arithmetic process.
// process Process to solve.
func solveProcess(process valueProcess) objects.Value {
	// Solve process.
	// operator Operator of process.
	// first First value.
	// second Second value.
	solve := func(operator objects.Token, first, second float64) float64 {
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
		} else {
			fract.Error(operator, "Operator is invalid!")
		}

		return result
	}

	value := objects.Value{}

	// String?
	if process.FirstV.Type == fract.VALString ||
		process.SecondV.Type == fract.VALString {
		value.Type = fract.VALString
		if process.FirstV.Type == process.SecondV.Type { // Both string?
			value.Content = []string{process.FirstV.Content[0] +
				process.SecondV.Content[0]}
			return value
		}

		if process.FirstV.Type == fract.VALString {
			if process.SecondV.Array {
				if len(process.SecondV.Content) == 0 {
					fract.Error(process.Second, "Array is empty!")
				}
				if len(process.FirstV.Content[0]) != len(process.SecondV.Content) &&
					(len(process.FirstV.Content[0]) != 1 && len(process.SecondV.Content) != 1) {
					fract.Error(process.Second,
						"Array element count is not one or equals to first array!")
				}
				var sb strings.Builder
				for _, char := range process.FirstV.Content[0] {
					if strings.Contains(process.SecondV.Content[0], grammar.TokenDot) {
						fract.Error(process.Second, "Float values cannot concatenate string values!")
					}
					sb.WriteRune(char + rune(arithmetic.ToArithmetic(process.SecondV.Content[0])))
				}
				value.Content = []string{sb.String()}
			} else {
				if process.SecondV.Type == fract.VALFloat {
					fract.Error(process.Second, "Float values cannot concatenate string values!")
				}
				var sb strings.Builder
				val := rune(arithmetic.ToArithmetic(process.SecondV.Content[0]))
				for _, char := range process.FirstV.Content[0] {
					sb.WriteRune(char + val)
				}
				value.Content = []string{sb.String()}
			}
		} else {
			if process.FirstV.Array {
				if len(process.FirstV.Content) == 0 {
					fract.Error(process.First, "Array is empty!")
				}
				if len(process.FirstV.Content[0]) != len(process.SecondV.Content) &&
					(len(process.FirstV.Content[0]) != 1 && len(process.SecondV.Content) != 1) {
					fract.Error(process.Second,
						"Array element count is not one or equals to first array!")
				}
				var sb strings.Builder
				for _, char := range process.SecondV.Content[0] {
					if strings.Contains(process.FirstV.Content[0], grammar.TokenDot) {
						fract.Error(process.Second, "Float values cannot concatenate string values!")
					}
					sb.WriteRune(char + rune(arithmetic.ToArithmetic(process.FirstV.Content[0])))
				}
				value.Content = []string{sb.String()}
			} else {
				if process.FirstV.Type == fract.VALFloat {
					fract.Error(process.First, "Float values cannot concatenate string values!")
				}
				var sb strings.Builder
				val := rune(arithmetic.ToArithmetic(process.FirstV.Content[0]))
				for _, char := range process.SecondV.Content[0] {
					sb.WriteRune(char + val)
				}
				value.Content = []string{sb.String()}
			}
		}
		return value
	}

	if process.FirstV.Array && process.SecondV.Array {
		if len(process.FirstV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		} else if len(process.SecondV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		}
		if len(process.FirstV.Content) != len(process.SecondV.Content) &&
			(len(process.FirstV.Content) != 1 && len(process.SecondV.Content) != 1) {
			fract.Error(process.Second,
				"Array element count is not one or equals to first array!")
		}

		if len(process.FirstV.Content) == 1 {
			first := arithmetic.ToArithmetic(process.FirstV.Content[0])
			for index, current := range process.SecondV.Content {
				process.SecondV.Content[index] = fmt.Sprintf("%g",
					solve(process.Operator, first, arithmetic.ToArithmetic(current)))
			}
			value.Content = process.SecondV.Content
		} else if len(process.SecondV.Content) == 1 {
			second := arithmetic.ToArithmetic(process.SecondV.Content[0])
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = fmt.Sprintf("%g",
					solve(process.Operator, arithmetic.ToArithmetic(current), second))
			}
			value.Content = process.FirstV.Content
		} else {
			for index, current := range process.FirstV.Content {
				process.FirstV.Content[index] = fmt.Sprintf("%g",
					solve(process.Operator, arithmetic.ToArithmetic(current),
						arithmetic.ToArithmetic(process.SecondV.Content[index])))
			}
			value.Content = process.FirstV.Content
		}
		value.Array = true
	} else if process.FirstV.Array {
		if len(process.FirstV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		}

		second := arithmetic.ToArithmetic(process.SecondV.Content[0])
		for index, current := range process.FirstV.Content {
			process.FirstV.Content[index] = fmt.Sprintf("%g",
				solve(process.Operator, arithmetic.ToArithmetic(current), second))
		}
		value.Array = true
		value.Content = process.FirstV.Content
	} else if process.SecondV.Array {
		if len(process.SecondV.Content) == 0 {
			fract.Error(process.First, "Array is empty!")
		}

		first := arithmetic.ToArithmetic(process.FirstV.Content[0])
		for index, current := range process.SecondV.Content {
			process.SecondV.Content[index] = fmt.Sprintf("%g",
				solve(process.Operator, arithmetic.ToArithmetic(current), first))
		}
		value.Array = true
		value.Content = process.SecondV.Content
	} else {
		value.Content = []string{fmt.Sprintf("%g",
			solve(process.Operator, arithmetic.ToArithmetic(process.FirstV.Content[0]),
				arithmetic.ToArithmetic(process.SecondV.Content[0])))}
	}

	return value
}

// __processValue Process value.
// first This is first value.
// token Token to process.
// tokens Tokens to process.
// index Index of token.
func (i *Interpreter) _processValue(first bool, operation *valueProcess,
	tokens *vector.Vector, index int) int {
	token := operation.First
	if !first {
		token = operation.Second
	}

	if token.Type == fract.TypeName {
		if index < len(tokens.Vals)-1 {
			next := tokens.Vals[index+1].(objects.Token)
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
					for ; cindex < len(tokens.Vals); cindex++ {
						current := tokens.Vals[cindex].(objects.Token)
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

					valueList := tokens.Sublist(index+2, cindex-index-3)
					// Index value is empty?
					if valueList.Vals == nil {
						fract.Error(token, "Index is not defined!")
					}

					value := i.processValue(valueList)
					if value.Array {
						fract.Error(tokens.Vals[index].(objects.Token),
							"Arrays is not used in index access!")
					} else if value.Type != fract.VALInteger {
						fract.Error(tokens.Vals[index].(objects.Token),
							"Only integer values can used in index access!")
					}
					position, err := arithmetic.ToInt(value.Content[0])
					if err != nil {
						fract.Error(tokens.Vals[index].(objects.Token),
							"Value out of range!")
					}

					variable := i.vars[vindex]

					if !variable.Value.Array && variable.Value.Type != fract.VALString {
						fract.Error(tokens.Vals[index].(objects.Token),
							"Index accessor is can not used with non-array variables!")
					}

					if variable.Value.Array {
						position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
					} else {
						position = parser.ProcessArrayIndex(len(variable.Value.Content[0]), position)
					}

					if position == -1 {
						fract.Error(tokens.Vals[index].(objects.Token),
							"Index is out of range!")
					}
					tokens.RemoveRange(index+1, cindex-index-1)

					var val string
					if variable.Value.Array {
						val = variable.Value.Content[position]
					} else {
						val = arithmetic.IntToString(variable.Value.Content[0][position])
					}

					if first {
						operation.FirstV.Content = []string{val}
						operation.FirstV.Array = false
						if variable.Value.Type == fract.VALString {
							operation.FirstV.Type = fract.VALString
						}
					} else {
						operation.SecondV.Content = []string{val}
						operation.SecondV.Array = false
						if variable.Value.Type == fract.VALString {
							operation.SecondV.Type = fract.VALString
						}
					}
					return 0
				} else if next.Value == grammar.TokenLParenthes { // Function?
					// Find close parentheses.
					cindex := index + 1
					bracketCount := 1
					for ; cindex < len(tokens.Vals); cindex++ {
						current := tokens.Vals[cindex].(objects.Token)
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
					value := i.processFunctionCall(*tokens.Sublist(index, cindex-index))
					if !operation.FirstV.Array && value.Content == nil {
						fract.Error(token, "Function is not return any value!")
					}
					tokens.RemoveRange(index+1, cindex-index-1)
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

		variable := i.vars[vindex]

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
				current := tokens.Vals[oindex].(objects.Token)
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
						tokens.Sublist(oindex, index-oindex+1)).Content
				} else {
					operation.SecondV.Array = true
					operation.SecondV.Content = i.processArrayValue(
						tokens.Sublist(oindex, index-oindex+1)).Content
				}
				tokens.RemoveRange(oindex, index-oindex)
				return index - oindex
			}

			endToken := tokens.Vals[oindex-1].(objects.Token)
			vindex := i.varIndexByName(endToken.Value)
			if vindex == -1 {
				fract.Error(endToken, "Variable is not defined in this name!: "+endToken.Value)
			}
			valueList := tokens.Sublist(oindex+1, index-oindex-1)
			// Index value is empty?
			if valueList.Vals == nil {
				fract.Error(endToken, "Index is not defined!")
			}

			value := i.processValue(valueList)
			if value.Array {
				fract.Error(tokens.Vals[index].(objects.Token),
					"Arrays is not used in index access!")
			} else if value.Type != fract.VALInteger {
				fract.Error(tokens.Vals[index].(objects.Token),
					"Only integer values can used in index access!")
			}

			position, err := arithmetic.ToInt(value.Content[0])
			if err != nil {
				fract.Error(tokens.Vals[oindex].(objects.Token),
					"Value out of range!")
			}

			variable := i.vars[vindex]

			if !variable.Value.Array && variable.Value.Type != fract.VALString {
				fract.Error(tokens.Vals[oindex].(objects.Token),
					"Index accessor is can not used with non-array variables!")
			}

			if variable.Value.Array {
				position = parser.ProcessArrayIndex(len(variable.Value.Content), position)
			} else {
				position = parser.ProcessArrayIndex(len(variable.Value.Content[0]), position)
			}

			if position == -1 {
				fract.Error(tokens.Vals[oindex].(objects.Token),
					"Index is out of range!")
			}
			tokens.RemoveRange(oindex-1, index-oindex+1)

			var val string
			if variable.Value.Array {
				val = variable.Value.Content[position]
			} else {
				val = arithmetic.IntToString(variable.Value.Content[0][position])
			}

			if first {
				operation.FirstV.Content = []string{val}
				operation.FirstV.Array = false
				if variable.Value.Type == fract.VALString {
					operation.FirstV.Type = fract.VALString
				}
			} else {
				operation.SecondV.Content = []string{val}
				operation.FirstV.Array = false
			}

			return index - oindex + 1
		} else if token.Value == grammar.TokenLBracket {
			// Array constructor.
			cindex := index + 1
			bracketCount := 1
			for ; cindex < len(tokens.Vals); cindex++ {
				current := tokens.Vals[cindex].(objects.Token)
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

			if first {
				operation.FirstV.Array = true
				operation.FirstV.Content = i.processArrayValue(
					tokens.Sublist(index, cindex-index+1)).Content
			} else {
				operation.SecondV.Array = true
				operation.SecondV.Content = i.processArrayValue(
					tokens.Sublist(index, cindex-index+1)).Content
			}
			tokens.RemoveRange(index+1, cindex-index-1)
			return 0
		} else if token.Value == grammar.TokenLBrace {
			// Array initializer.

			// Find close brace.
			cindex := index + 1
			braceCount := 1
			for ; cindex < len(tokens.Vals); cindex++ {
				current := tokens.Vals[cindex].(objects.Token)
				if current.Type == fract.TypeBrace {
					if current.Value == grammar.TokenLBrace {
						fract.Error(current, "Arrays is can not take array value as element!")
						braceCount++
					} else if current.Value == grammar.TokenRBrace {
						braceCount--
						if braceCount == 0 {
							break
						}
					}
				}
			}

			value := i.processArrayValue(tokens.Sublist(index, cindex-index+1))
			if first {
				operation.FirstV = value
			} else {
				operation.SecondV = value
			}
			tokens.RemoveRange(index+1, cindex-index-1)
			return 0
		} else if token.Value == grammar.TokenRBrace {
			// Array initializer.

			// Find open brace.
			braceCount := 1
			oindex := index - 1
			nestedArray := false
			for ; oindex >= 0; oindex-- {
				current := tokens.Vals[oindex].(objects.Token)
				if current.Type == fract.TypeBrace {
					if current.Value == grammar.TokenRBrace {
						braceCount++
						nestedArray = true
					} else if current.Value == grammar.TokenLBrace {
						if nestedArray {
							fract.Error(current, "Arrays is can not take array value as element!")
						}
						braceCount--
						if braceCount == 0 {
							break
						}
					}
				}
			}

			value := i.processArrayValue(tokens.Sublist(oindex, index-oindex+1))
			if first {
				operation.FirstV = value
			} else {
				operation.SecondV = value
			}
			tokens.RemoveRange(oindex, index-oindex)
			return index - oindex
		} else if token.Value == grammar.TokenRParenthes {
			// Function.

			// Find open parentheses.
			bracketCount := 1
			oindex := index - 1
			for ; oindex >= 0; oindex-- {
				current := tokens.Vals[oindex].(objects.Token)
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
			value := i.processFunctionCall(*tokens.Sublist(oindex, index-oindex+1))
			if value.Content == nil {
				fract.Error(tokens.Vals[oindex].(objects.Token),
					"Function is not return any value!")
			}
			if first {
				operation.FirstV = value
			} else {
				operation.SecondV = value
			}
			tokens.RemoveRange(oindex, index-oindex)
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
		_, err := arithmetic.ToFloat64(token.Value)
		if err != nil {
			fract.Error(token, "Value out of range!")
		}
	}

	if first {
		operation.FirstV.Array = false
		if strings.HasPrefix(token.Value, grammar.TokenQuote) ||
			strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) { // String?
			operation.FirstV.Type = fract.VALString
			operation.FirstV.Content = []string{token.Value[1 : len(token.Value)-1]}
		} else {
			operation.FirstV.Content = []string{token.Value}
		}
	} else {
		operation.SecondV.Array = false
		if strings.HasPrefix(token.Value, grammar.TokenQuote) ||
			strings.HasPrefix(token.Value, grammar.TokenDoubleQuote) { // String?
			operation.SecondV.Type = fract.VALString
			operation.SecondV.Content = []string{token.Value[1 : len(token.Value)-1]}
		} else {
			operation.SecondV.Content = []string{token.Value}
		}
	}

	// Type check.
	if token.Type == fract.TypeBooleanTrue ||
		token.Type == fract.TypeBooleanFalse {
		if first {
			operation.FirstV.Type = fract.VALBoolean
		} else {
			operation.SecondV.Type = fract.VALBoolean
		}
	} else if strings.Contains(token.Value, grammar.TokenDot) { // Float?
		if first {
			operation.FirstV.Type = fract.VALFloat
		} else {
			operation.SecondV.Type = fract.VALFloat
		}
	}

	return 0
}

// processArrayValue Process array value.
// tokens Tokens.
func (i *Interpreter) processArrayValue(tokens *vector.Vector) objects.Value {
	value := objects.Value{
		Array: true,
	}

	first := tokens.Vals[0].(objects.Token)

	// Initializer?
	if first.Value == grammar.TokenLBracket {
		valueList := tokens.Sublist(1, len(tokens.Vals)-2)

		if valueList.Vals == nil {
			fract.Error(first, "Size is not defined!")
		}

		value := i.processValue(valueList)
		if value.Array {
			fract.Error(first, "Arrays is not used in array constructors!")
		} else if value.Type != fract.VALInteger {
			fract.Error(first, "Only integer values can used in array constructors!")
		}

		val, _ := arithmetic.ToInt64(value.Content[0])
		if val < 0 {
			fract.Error(first, "Value is not lower than zero!")
		}
		value.Content = make([]string, val)
		for index := range value.Content {
			value.Content[index] = "0"
		}
		return value
	}

	comma := 1
	for index := 1; index < len(tokens.Vals)-1; index++ {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeComma {
			lst := tokens.Sublist(comma, index-comma)
			if lst.Vals == nil {
				fract.Error(first, "Value is not defined!")
			}
			val := i.processValue(lst)
			value.Content = append(value.Content, val.Content...)
			if value.Type != fract.VALString {
				value.Type = val.Type
			}
			comma = index + 1
		}
	}

	if comma < len(tokens.Vals)-1 {
		lst := tokens.Sublist(comma, len(tokens.Vals)-comma-1)
		if lst.Vals == nil {
			fract.Error(first, "Value is not defined!")
		}
		val := i.processValue(lst)
		value.Content = append(value.Content, val.Content...)
		if value.Type != fract.VALString {
			value.Type = val.Type
		}
	}

	return value
}

// processValue Process value.
// tokens Tokens to process.
func (i *Interpreter) processValue(tokens *vector.Vector) objects.Value {
	value := objects.Value{
		Content: []string{"0"},
		Array:   false,
	}

	i.processRange(tokens)

	// Is conditional expression?
	brace := 0
	for _, current := range tokens.Vals {
		current := current.(objects.Token)
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
			value.Content = []string{i.processCondition(tokens)}
			return value
		}
	}

	// Calculate data count.
	data_count := 0
	bracket := 0
	for _, current := range tokens.Vals {
		current := current.(objects.Token)
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLBracket ||
				current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLParenthes {
				bracket++
			} else {
				bracket--
			}
		}
		if bracket > 0 {
			continue
		}
		if current.Type == fract.TypeValue ||
			current.Type == fract.TypeName {
			data_count++
		}
	}
	data_count -= 1

	// Decompose arithmetic operations.
	priorityIndex := parser.IndexProcessPriority(*tokens)
	looped := priorityIndex != -1
	for priorityIndex != -1 {
		data_count--
		var operation valueProcess

		operation.First = tokens.Vals[priorityIndex-1].(objects.Token)
		priorityIndex -= i._processValue(true, &operation,
			tokens, priorityIndex-1)
		operation.Operator = tokens.Vals[priorityIndex].(objects.Token)

		operation.Second = tokens.Vals[priorityIndex+1].(objects.Token)
		priorityIndex -= i._processValue(false, &operation,
			tokens, priorityIndex+1)

		resultValue := solveProcess(operation)

		operation.Operator.Value = grammar.TokenPlus
		operation.Second = tokens.Vals[priorityIndex+1].(objects.Token)
		operation.FirstV = value
		operation.SecondV = resultValue

		resultValue = solveProcess(operation)
		value = resultValue

		// Remove processed processes.
		tokens.RemoveRange(priorityIndex-1, 3)
		tokens.Insert(priorityIndex-1, objects.Token{Value: "0"})

		// Find next operator.
		priorityIndex = parser.IndexProcessPriority(*tokens)
	}

	// Not operatored?
	if !looped {
		var operation valueProcess
		operation.First = tokens.Vals[0].(objects.Token)
		operation.FirstV.Array = true // Ignore nil control if function call.
		i._processValue(true, &operation, tokens, 0)
		value = operation.FirstV
	}

	if data_count > 0 {
		fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token),
			"Invalid value!")
	}

	return value
}
