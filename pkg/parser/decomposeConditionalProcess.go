/*
	DecomposeConditionalProcess Function.
*/

package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/grammar"
	obj "github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

// findNextOrOperator Find next or condition operator index and return if find, return -1 if not.
// tokens Tokens to find.
// pos Position of start to find.
// operator Operator to find.
func findNextOperator(tokens []obj.Token, pos int, operator string) int {
	brace := 0
	for ; pos < len(tokens); pos++ {
		current := tokens[pos]
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLBrace ||
				current.Value == grammar.TokenLBracket ||
				current.Value == grammar.TokenLParenthes {
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

// DecomposeConditionalProcess Decompose and returns conditional expressions by operators.
// tokens Tokens to process.
// operator Operator to decompose.
func DecomposeConditionalProcess(tokens []obj.Token, operator string) *[][]obj.Token {
	var expressions [][]obj.Token

	last := 0
	index := findNextOperator(tokens, last, operator)
	if index == 0 { // Operator is first element of vector?
		fract.Error(tokens[0], "Operator spam!")
	}
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
