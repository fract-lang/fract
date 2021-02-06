/*
	DecomposeOrConditionalProcess Function.
*/

package parser

import (
	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// findNextOrOperator Find next or condition operator index and return if find, return -1 if not.
// tokens Tokens to find.
// pos Position of start to find.
func findNextOrOperator(tokens *vector.Vector, pos int) int {
	for ; pos < tokens.Len(); pos++ {
		current := tokens.At(pos).(objects.Token)
		if current.Type == fract.TypeOperator && current.Value == grammar.TokenVerticalBar {
			return pos
		}
	}
	return -1
}

// DecomposeOrConditionalProcess Decompose and returns conditional expressions by operators.
// tokens Tokens to process.
func DecomposeOrConditionalProcess(tokens *vector.Vector) vector.Vector {
	expressions := *vector.New()

	last := 0
	index := findNextOrOperator(tokens, last)
	if index == 0 { // Operator is first element of vector?
		fract.Error(tokens.First().(objects.Token), "Operator spam!")
	}
	for index != -1 {
		expressions.Append(tokens.Sublist(last, index-last))
		last = index + 1
		index = findNextOrOperator(tokens, last) // Find next.
		if index == tokens.Len()-1 {
			fract.Error(tokens.Last().(objects.Token), "Operator defined, but for what?")
		}
	}
	if last != tokens.Len() {
		expressions.Append(tokens.Sublist(last, tokens.Len()-last))
	}

	return expressions
}
