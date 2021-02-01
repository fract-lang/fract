package parser

import (
	"strings"

	"../fract"
	"../grammar"
	"../objects"
	"../utilities/vector"
)

// DecomposeArithmeticProcesses Decompose and returns arithmetic processes by operators.
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
