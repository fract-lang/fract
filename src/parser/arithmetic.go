/*
	ARITHMETIC FUNCTIONS
*/

package parser

import (
	"strings"

	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// DecomposeArithmeticProcesses Decompose and returns arithmetic processes by operators.
// tokens Tokens to process.
func DecomposeArithmeticProcesses(tokens *vector.Vector) vector.Vector {
	var (
		operator bool
		last     objects.Token
	)
	processes := *vector.New()

	for index := 0; index < tokens.Len(); index++ {
		_token := tokens.At(index).(objects.Token)
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
			operator = index < tokens.Len()-1
		} else {
			fract.Error(_token, "Invalid value!")
		}
	}

	if last.Type == fract.TypeOperator {
		fract.Error(processes.Last().(objects.Token), "Operator defined, but for what?")
	}

	return processes
}

// IndexProcessPriority Find index of priority operator.
// Returns index of operator if found, returns -1 if not.
//
// tokens Tokens to search.
func IndexProcessPriority(tokens *vector.Vector) int {
	// Returns -1 if vector contains one value.
	if len(tokens.Vals) == 1 {
		return -1
	}

	/* Find exponentiation. */
	for index := 0; index < len(tokens.Vals); index++ {
		if tokens.At(index).(objects.Token).Value == grammar.TokenCaret {
			return index
		}
	}

	/* Find mod. */
	for index := 0; index < len(tokens.Vals); index++ {
		if tokens.At(index).(objects.Token).Value == grammar.TokenPercent {
			return index
		}
	}

	/* Find multipy or divide. */
	for index := 0; index < len(tokens.Vals); index++ {
		_token := tokens.At(index).(objects.Token)
		if _token.Value == grammar.TokenStar ||
			_token.Value == grammar.TokenSlash ||
			_token.Value == grammar.TokenReverseSlash ||
			_token.Value == grammar.IntegerDivision ||
			_token.Value == grammar.IntegerDivideWithBigger {
			return index
		}
	}

	/* Addition or subtraction. */
	/*for index := 0; index < len(tokens.Vals); index++ {
		_token := tokens.At(index).(objects.Token)
		if _token.Value == grammar.TokenPlus || _token.Value == grammar.TokenMinus {
			return index
		}
	}*/

	return 1
}
