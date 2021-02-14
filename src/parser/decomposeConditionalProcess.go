/*
	DecomposeConditionalProcess Function.
*/

package parser

import (
	"../fract"
	"../objects"
	"../utilities/vector"
)

// findNextOrOperator Find next or condition operator index and return if find, return -1 if not.
// tokens Tokens to find.
// pos Position of start to find.
// operator Operator to find.
func findNextOperator(tokens *vector.Vector, pos int, operator string) int {
	for ; pos < len(tokens.Vals); pos++ {
		current := tokens.Vals[pos].(objects.Token)
		if current.Type == fract.TypeOperator && current.Value == operator {
			return pos
		}
	}
	return -1
}

// DecomposeConditionalProcess Decompose and returns conditional expressions by operators.
// tokens Tokens to process.
// operator Operator to decompose.
func DecomposeConditionalProcess(tokens *vector.Vector, operator string) vector.Vector {
	expressions := *vector.New()

	last := 0
	index := findNextOperator(tokens, last, operator)
	if index == 0 { // Operator is first element of vector?
		fract.Error(tokens.Vals[0].(objects.Token), "Operator spam!")
	}
	for index != -1 {
		expressions.Vals = append(expressions.Vals, tokens.Sublist(last, index-last))
		last = index + 1
		index = findNextOperator(tokens, last, operator) // Find next.
		if index == len(tokens.Vals)-1 {
			fract.Error(tokens.Vals[len(tokens.Vals)-1].(objects.Token),
				"Operator defined, but for what?")
		}
	}
	if last != len(tokens.Vals) {
		expressions.Vals = append(expressions.Vals, tokens.Sublist(last, len(tokens.Vals)-last))
	}

	return expressions
}
