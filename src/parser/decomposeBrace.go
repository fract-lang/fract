/*
	DecomposeBrace Function
*/

package parser

import (
	"github.com/fract-lang/fract/src/fract"
	"github.com/fract-lang/fract/src/objects"
	"github.com/fract-lang/fract/src/utils/vector"
)

// DecomposeBrace Returns range tokens and index of first parentheses.
// Remove range tokens from original tokens.
//
// tokens Tokens to process.
// open Open bracket.
// close Close bracket.
// nonCheck Check empty bracket content.
func DecomposeBrace(tokens *vector.Vector, open, close string,
	nonCheck bool) (vector.Vector, int) {
	var (
		first int = -1
		last  int
	)

	/* Find open parentheses. */
	if nonCheck {
		name := false
		for index := range tokens.Vals {
			current := tokens.Vals[index].(objects.Token)
			if current.Type == fract.TypeName {
				name = true
			} else if !name && current.Type == fract.TypeBrace && current.Value == open {
				first = index
				break
			}
		}
	} else {
		for index := range tokens.Vals {
			current := tokens.Vals[index].(objects.Token)
			if current.Type == fract.TypeBrace && current.Value == open {
				first = index
				break
			}
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
	_range := *tokens.Sublist(first+1, length)

	// Bracket content is empty?
	if nonCheck && _range.Vals == nil {
		fract.Error(tokens.Vals[first].(objects.Token), "Brackets content are empty!")
	}

	/* Remove range from original tokens. */
	tokens.RemoveRange(first, last-first+1)

	return _range, first
}
