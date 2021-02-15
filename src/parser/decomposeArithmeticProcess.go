/*
	DecomposeArithmeticProcess Function
*/

package parser

import (
	"fmt"
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

	for index := range tokens.Vals {
		_token := tokens.Vals[index].(objects.Token)
		if _token.Type == fract.TypeOperator {
			if !operator {
				fract.Error(_token, "Operator spam!")
			}
			last = _token
			processes.Vals = append(processes.Vals, _token)
			operator = false
		} else if _token.Type == fract.TypeValue || _token.Type == fract.TypeName ||
			_token.Type == fract.TypeBooleanTrue || _token.Type == fract.TypeBooleanFalse ||
			(_token.Type == fract.TypeBrace && (_token.Value == grammar.TokenLBracket ||
				_token.Value == grammar.TokenRBracket)) {
			if _token.Type == fract.TypeName && last.Type == fract.TypeOperator &&
				last.Value == grammar.TokenMinus &&
				strings.HasPrefix(_token.Value, grammar.TokenMinus) {
				fract.Error(_token, "Negative number declare after subtraction!")
			}
			last = _token
			processes.Vals = append(processes.Vals, _token)
			operator = index < len(tokens.Vals)-1
		} else {
			fract.Error(_token, "Invalid value!")
		}
	}

	if last.Type == fract.TypeOperator {
		fract.Error(processes.Vals[len(processes.Vals)-1].(objects.Token),
			"Operator defined, but for what?")
	}

	fmt.Println(processes.Vals, "\n")
	return processes
}
