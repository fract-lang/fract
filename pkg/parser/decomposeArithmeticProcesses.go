package parser

import (
	"github.com/fract-lang/fract/pkg/fract"
	"github.com/fract-lang/fract/pkg/objects"
)

// DecomposeArithmeticProcesses returns arithmetic processes by operators.
func DecomposeArithmeticProcesses(tokens []objects.Token) *[]objects.Token {
	var (
		operator  bool
		last      objects.Token
		processes []objects.Token
	)

	for index, token := range tokens {
		if token.Type == fract.TypeOperator {
			if !operator {
				fract.Error(token, "Operator spam!")
			}
			last = token
			processes = append(processes, token)
			operator = false
		} else if token.Type == fract.TypeValue || token.Type == fract.TypeName ||
			token.Type == fract.TypeBooleanTrue || token.Type == fract.TypeBooleanFalse ||
			token.Type == fract.TypeBrace || token.Type == fract.TypeComma {
			last = token
			processes = append(processes, token)
			operator = index < len(tokens)-1
		} else {
			fract.Error(token, "Invalid value!")
		}
	}

	if last.Type == fract.TypeOperator {
		fract.Error(processes[len(processes)-1], "Operator defined, but for what?")
	}

	return &processes
}
