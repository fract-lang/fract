package parser

import (
	"strings"

	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// DecomposeBrace Returns range tokens and index of first parentheses.
// Remove range tokens from original tokens.
//
// tokens Tokens to process.
func DecomposeBrace(tokens *vector.Vector) (vector.Vector, int) {
	var (
		first int = -1
		last  int
	)

	len := len(tokens.Vals)

	/* Find open parentheses. */
	for index := 0; index < len; index++ {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeBrace {
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
	for index := first + 1; index < len; index++ {
		current := tokens.Vals[index].(objects.Token)
		if current.Type == fract.TypeBrace {
			if current.Value == grammar.TokenLParenthes {
				count++
			} else {
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

	/* Remove range from original tokens. */
	tokens.RemoveRange(first, last-first+1)

	return _range, first
}

// DecomposeArithmeticProcesses Decompose and returns arithmetic processes by operators.
// tokens Tokens to process.
func DecomposeArithmeticProcesses(tokens *vector.Vector) vector.Vector {
	var (
		operator bool
		last     objects.Token
	)
	processes := *vector.New()
	len := len(tokens.Vals)

	for index := 0; index < len; index++ {
		_token := tokens.Vals[index].(objects.Token)
		if _token.Type == fract.TypeOperator {
			if !operator {
				fract.Error(_token, "Operator spam!")
			}
			last = _token
			processes.Append(_token)
			operator = false
		} else if _token.Type == fract.TypeValue {
			if last.Type == fract.TypeOperator && last.Value == grammar.TokenMinus &&
				strings.HasPrefix(_token.Value, grammar.TokenMinus) {
				fract.Error(_token, "Negative number declare after subtraction!")
			}
			last = _token
			processes.Append(_token)
			operator = index < len-1
		} else {
			fract.Error(_token, "Invalid value!")
		}
	}

	if last.Type == fract.TypeOperator {
		fract.Error(processes.Last().(objects.Token), "Operator defined, but for what?")
	}

	return processes
}
