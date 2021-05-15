package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
	"github.com/fract-lang/fract/pkg/vector"
)

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

	/*
		Skip find close parentheses and result ready steps
		if open parentheses is not found.
	*/
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

			if count == 0 { break }
		}
		length++
	}

	_range := vector.Sublist(*tokens, first+1, length)

	// Bracket content is empty?
	if nonCheck && _range == nil {
		fract.Error((*tokens)[first], "Brackets content are empty!")
	}

	/* Remove range from original tokens. */
	vector.RemoveRange(tokens, first, (first + length + 1)-first+1)

	if _range == nil {
		return nil, first
	}

	return *_range, first
}
