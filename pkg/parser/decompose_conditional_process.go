package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

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
