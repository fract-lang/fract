package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// Check arithmetic processes validity.
func CheckArithmeticProcesses(tokens []objects.Token) {
	var (
		operator bool
		brace    int
	)
	for index := 0; index < len(tokens); index++ {
		switch token := tokens[index]; token.Type {
		case fract.TypeOperator:
			if !operator {
				fract.Error(token, "Operator spam!")
			}
			operator = false
		case fract.TypeValue, fract.TypeName, fract.TypeComma, fract.TypeBrace:
			switch token.Type {
			case fract.TypeBrace:
				if token.Value == "(" || token.Value == "[" || token.Value == "{" {
					brace++
				} else {
					brace--
				}
			case fract.TypeComma:
				if brace == 0 {
					fract.Error(token, "Invalid syntax!")
				}
			}
			operator = index < len(tokens)-1
		default:
			fract.Error(token, "Invalid syntax!")
		}
	}
	if tokens[len(tokens)-1].Type == fract.TypeOperator {
		fract.Error(tokens[len(tokens)-1], "Operator overflow!")
	}
}

// DecomposeBrace returns range tokens and index of first parentheses.
// Remove range tokens from original tokens.
func DecomposeBrace(tokens *[]objects.Token, open, close string, nonCheck bool) ([]objects.Token, int) {
	first := -1
	/* Find open parentheses. */
	if nonCheck {
		name := false
		for index, current := range *tokens {
			if current.Type == fract.TypeName {
				name = true
			} else if !name && current.Type == fract.TypeBrace && current.Value == open {
				first = index
				break
			} else {
				name = false
			}
		}
	} else {
		for index, current := range *tokens {
			if current.Type == fract.TypeBrace && current.Value == open {
				first = index
				break
			}
		}
	}
	// Skip find close parentheses and result ready steps
	// if open parentheses is not found.
	if first == -1 {
		return nil, -1
	}
	// Find close parentheses.
	count := 1
	length := 0
	for index := first + 1; index < len(*tokens); index++ {
		current := (*tokens)[index]
		if current.Type == fract.TypeBrace {
			if current.Value == open {
				count++
			} else if current.Value == close {
				count--
			}
			if count == 0 {
				break
			}
		}
		length++
	}

	_range := vector.Sublist(*tokens, first+1, length)
	// Bracket content is empty?
	if nonCheck && _range == nil {
		fract.Error((*tokens)[first], "Brackets content are empty!")
	}
	/* Remove range from original tokens. */
	vector.RemoveRange(tokens, first, (first+length+1)-first+1)
	if _range == nil {
		return nil, first
	}
	return *_range, first
}

// ProcessArrayIndex process array index by length.
func ProcessArrayIndex(length, index int) int {
	if index >= 0 {
		if index >= length {
			return -1
		}
		return index
	}
	index = length + index
	if index < 0 || index >= length {
		return -1
	}
	return index
}

// IsBlockStatement returns true if tokens is block start, return false if not.
func IsBlockStatement(tokens []objects.Token) bool {
	if tokens[0].Type == fract.TypeMacro { // Remove macro token.
		tokens = tokens[1:]
	}

	switch tokens[0].Type {
	case fract.TypeIf,
		fract.TypeLoop,
		fract.TypeFunction,
		fract.TypeTry:
		return true
	case fract.TypeProtected:
		if len(tokens) > 1 {
			if tokens[1].Type == fract.TypeFunction {
				return true
			}
		}
	}
	return false
}

// IndexProcessPriority find index of priority operator and
// returns index of operator if found, returns -1 if not.
func IndexProcessPriority(tokens []objects.Token) int {
	bracket := 0
	multiplyOrDivive := -1
	binaryOrAnd := -1
	additionOrSubtraction := -1
	for index, token := range tokens {
		if token.Type == fract.TypeBrace {
			if token.Value == "[" || token.Value == "{" || token.Value == "(" {
				bracket++
			} else {
				bracket--
			}
		}
		if bracket > 0 {
			continue
		}
		// Exponentiation or shifts.
		if token.Value == grammar.LeftBinaryShift || token.Value == grammar.RightBinaryShift ||
			token.Value == grammar.Exponentiation {
			return index
		}
		switch token.Value {
		case "%": // Modulus.
			return index
		case "*", "/", "\\", grammar.IntegerDivision, grammar.IntegerDivideWithBigger: // Multiply or division.
			if multiplyOrDivive == -1 {
				multiplyOrDivive = index
			}
		case "+", "-": // Addition or subtraction.
			if additionOrSubtraction == -1 {
				additionOrSubtraction = index
			}
		case "&", "|":
			if binaryOrAnd == -1 {
				binaryOrAnd = index
			}
		}
	}

	if multiplyOrDivive != -1 {
		return multiplyOrDivive
	} else if binaryOrAnd != -1 {
		return binaryOrAnd
	} else if additionOrSubtraction != -1 {
		return additionOrSubtraction
	}
	return -1
}

// FindConditionOperator return next condition operator.
func FindConditionOperator(tokens []objects.Token) (int, objects.Token) {
	for index, current := range tokens {
		if (current.Type == fract.TypeOperator && (current.Value == grammar.Equals ||
			current.Value == grammar.NotEquals || current.Value == ">" || current.Value == "<" ||
			current.Value == grammar.GreaterEquals || current.Value == grammar.LessEquals)) ||
			current.Type == fract.TypeIn {
			return index, current
		}
	}
	var token objects.Token
	return -1, token
}

// findNextOrOperator find next or condition operator index and return if find, return -1 if not.
func findNextOperator(tokens []objects.Token, pos int, operator string) int {
	brace := 0
	for ; pos < len(tokens); pos++ {
		current := tokens[pos]
		if current.Type == fract.TypeBrace {
			if current.Value == "{" || current.Value == "[" || current.Value == "(" {
				brace++
			} else {
				brace--
			}
		}
		if brace > 0 {
			continue
		}
		if current.Type == fract.TypeOperator && current.Value == operator {
			return pos
		}
	}
	return -1
}

// DecomposeConditionalProcess returns conditional expressions by operators.
func DecomposeConditionalProcess(tokens []objects.Token, operator string) *[][]objects.Token {
	var expressions [][]objects.Token
	last := 0
	index := findNextOperator(tokens, last, operator)
	for index != -1 {
		if index-last == 0 {
			fract.Error(tokens[last], "Where is the condition?")
		}
		expressions = append(expressions, *vector.Sublist(tokens, last, index-last))
		last = index + 1
		index = findNextOperator(tokens, last, operator) // Find next.
		if index == len(tokens)-1 {
			fract.Error(tokens[len(tokens)-1], "Operator defined, but for what?")
		}
	}
	if last != len(tokens) {
		expressions = append(expressions, *vector.Sublist(tokens, last, len(tokens)-last))
	}
	return &expressions
}
