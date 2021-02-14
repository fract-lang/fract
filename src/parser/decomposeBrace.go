/*
	DecomposeBrace Function
*/

package parser

import (
	"../fract"
	"../objects"
	"../utilities/vector"
)

// DecomposeBrace Returns range tokens and index of first parentheses.
// Remove range tokens from original tokens.
//
// tokens Tokens to process.
// open Open bracket.
// close Close bracket.
func DecomposeBrace(tokens *vector.Vector, open string, close string) (vector.Vector, int) {
	var (
		first int = -1
		last  int
	)

	/* Find open parentheses. */
	for index := range tokens.Vals {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeBrace && current.Value == open {
			first = index
			break
		}
	}

	/*
		Skip find close parentheses and result ready steps
		if open parentheses is not found.
	*/
	if first == -1 {
		return *new(vector.Vector), -1
	}

	/* Find close parentheses. */
	count := 1
	length := 0
	for index := first + 1; index < len(tokens.Vals); index++ {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeBrace {
			if current.Value == open {
				count++
			} else if current.Value == close {
				count--
			}
			if count == 0 {
				last = index
				break
			}
		}
		length++
	}
	_range := tokens.Sublist(first+1, length)

	// Bracket content is empty?
	if len(_range.Vals) == 0 {
		fract.Error(tokens.Vals[first].(objects.Token), "Brackets content are empty!")
	}

	/* Remove range from original tokens. */
	tokens.RemoveRange(first, last-first+1)

	return _range, first
}
